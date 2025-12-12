package iam

import (
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
	"github.com/IllumiKnowLabs/labstore/backend/internal/security"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type User struct {
	UserID      int            `db:"user_id"`
	Name        string         `db:"name"`
	AccessKeyID sql.NullString `db:"access_key"`
	SecretKey   []byte         `db:"secret_key"`
	Salt        []byte         `db:"salt"`

	GroupIDs  []string
	PolicyIDs []string
}

type CreateUserResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ CreateUserResponse"`
	CreateUserResult *CreateUserResult
}

type CreateUserResult struct {
	User             *UserResult
	ResponseMetadata *ResponseMetadata
}

type UserResult struct {
	// XMLName  xml.Name `xml:"User"`
	Path     string
	UserName string
	UserId   string
	Arn      string
}

type CreateAccessKeyResponse struct {
	XMLName               xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ CreateAccessKeyResponse"`
	CreateAccessKeyResult *CreateAccessKeyResult
}

type CreateAccessKeyResult struct {
	AccessKey        *AccessKeyResult
	ResponseMetadata *ResponseMetadata
}

type AccessKeyResult struct {
	XMLName         xml.Name `xml:"AccessKey"`
	UserName        string
	AccessKeyId     string
	Status          AccessKeyStatus
	SecretAccessKey string
}

type AccessKeyStatus string

const (
	AccessKeyActive   AccessKeyStatus = "Active"
	AccessKeyInactive AccessKeyStatus = "Inactive"
	AccessKeyExpired  AccessKeyStatus = "Expired"
)

func GetUserByAccessKey(accessKey string) (*User, error) {
	if user, ok := store.Users[accessKey]; ok {
		return user, nil
	}

	var user User

	if err := store.readDB.Get(&user, `SELECT * FROM users WHERE access_key = $1`, accessKey); err != nil {
		return nil, err
	}

	store.Users[user.AccessKeyID.String] = &user

	return &user, nil
}

func GetUserByName(name string) (*User, error) {
	var user User

	if err := store.readDB.Get(&user, `SELECT * FROM users WHERE name = $1`, name); err != nil {
		return nil, err
	}

	if user.AccessKeyID.Valid {
		store.Users[user.AccessKeyID.String] = &user
	}

	return &user, nil
}

func CreateUser(name string) (*User, *errs.IAMError) {
	if name == config.Admin.Auth.AccessKey {
		return nil, errs.IAMEntityAlreadyExists(name)
	}

	user := &User{Name: name}

	_, err := store.writeDB.NamedExec(`INSERT INTO users (name) VALUES (:name)`, &user)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				slog.Error("create user insert", "err", sqliteErr)
				return nil, errs.IAMEntityAlreadyExists(name)
			}
		}

		slog.Error("create user insert", "err", err)
		return nil, errs.IAMServiceFailure()
	}

	user, err = GetUserByName(name)
	if err != nil {
		slog.Error("get user by name", "err", err)
		return nil, errs.IAMServiceFailure()
	}

	return user, nil
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("UserName")
	if userName == "" {
		errs.Handle(w, errors.New("missing query parameter: UserName"))
		return
	}

	user, err := CreateUser(userName)
	if err != nil {
		errs.Handle(w, err)
		return
	}

	userPath := "/"

	response := &CreateUserResponse{
		CreateUserResult: &CreateUserResult{
			User: &UserResult{
				Path:     userPath,
				UserName: user.Name,
				UserId:   fmt.Sprint(user.UserID),
				Arn:      fmt.Sprintf("arn:labstore:iam::%d:%s%s", user.UserID, userPath, user.Name),
			},
			ResponseMetadata: &ResponseMetadata{
				RequestId: core.NewRequestID(),
			},
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}

func CreateAccessKey(user *User) (string, error) {
	secretKey, err := security.GeneratePassword(42)
	if err != nil {
		return "", err
	}

	user.AccessKeyID = sql.NullString{
		String: user.Name,
		Valid:  user.Name != "",
	}

	encryptedSecretKey, err := security.EncryptAESGCM(secretKey, config.Storage.MasterKeyPath)
	if err != nil {
		return "", err
	}

	user.SecretKey = encryptedSecretKey.Value
	user.Salt = encryptedSecretKey.Salt

	query := `
	UPDATE users
	SET
		access_key = :access_key,
		secret_key = :secret_key,
		salt = :salt
	WHERE user_id = :user_id
	`

	if _, err := store.writeDB.NamedExec(query, user); err != nil {
		return "", err
	}

	return secretKey, nil
}

func CreateAccessKeyHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("UserName")
	if userName == "" {
		errs.Handle(w, errors.New("missing query parameter: UserName"))
		return
	}

	user, err := GetUserByName(userName)
	if err != nil {
		slog.Error("get user by name", "err", err)
		errs.Handle(w, fmt.Errorf("user not found: %s", userName))
		return
	}

	secretKey, err := CreateAccessKey(user)
	if err != nil {
		slog.Error("create access key", "err", err)
		errs.Handle(w, errors.New("could not create access key"))
		return
	}

	response := &CreateAccessKeyResponse{
		CreateAccessKeyResult: &CreateAccessKeyResult{
			AccessKey: &AccessKeyResult{
				UserName:        user.Name,
				AccessKeyId:     user.AccessKeyID.String,
				Status:          AccessKeyActive,
				SecretAccessKey: secretKey,
			},
			ResponseMetadata: &ResponseMetadata{
				RequestId: core.NewRequestID(),
			},
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}

func (user *User) EncryptedData() *security.EncryptedData {
	return &security.EncryptedData{
		Value: user.SecretKey,
		Salt:  user.Salt,
	}
}

func (store *Store) setupAdmin() error {
	encryptedData, err := security.EncryptAESGCM(config.Admin.Auth.SecretKey, config.Storage.MasterKeyPath)
	if err != nil {
		return err
	}

	store.Users[config.Admin.Auth.AccessKey] = &User{
		Name:        "Administrator",
		AccessKeyID: sql.NullString{String: config.Admin.Auth.AccessKey, Valid: true},
		SecretKey:   encryptedData.Value,
		Salt:        encryptedData.Salt,
		PolicyIDs:   []string{adminPolicy},
	}

	store.Policies[adminPolicy] = &Policy{
		ID: adminPolicy,
		Document: &PolicyDocument{
			Version: latestPolicyDocumentVersion,
			Statement: []Statement{
				{
					Effect:    allow,
					Actions:   []Action{Action(Any)},
					Resources: []string{Any},
				},
			},
		},
	}

	return nil
}
