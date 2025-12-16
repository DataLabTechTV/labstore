package iam

import (
	"database/sql"
	"encoding/xml"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
	"github.com/IllumiKnowLabs/labstore/backend/internal/security"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

const (
	defaultAdminUserName = "Administrator"
	defaultUserPath      = "/"
)

type User struct {
	UserID string `db:"user_id"`
	Name   string `db:"name"`
	Arn    string `db:"arn"`

	AccessKeyID sql.NullString `db:"access_key"`
	SecretKey   []byte         `db:"secret_key"`
	Salt        []byte         `db:"salt"`

	GroupIDs  []string
	PolicyIDs []string
}

type cachedUser struct {
	user        *User
	loadedAt    time.Time
	neverExpire bool
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

func (user *User) EncryptedData() *security.EncryptedData {
	return &security.EncryptedData{
		Value: user.SecretKey,
		Salt:  user.Salt,
	}
}

func (store *Store) GetUserByAccessKey(accessKey string) (*User, error) {
	if cachedUser, ok := store.Users[accessKey]; ok {
		if cachedUser.neverExpire || time.Since(cachedUser.loadedAt) < store.ttl {
			return cachedUser.user, nil
		}

		slog.Debug("invalidating cached user", "accessKey", accessKey)
		delete(store.Users, accessKey)
	}

	var user User
	if err := store.readDB.Get(&user, `SELECT * FROM users WHERE access_key = $1`, accessKey); err != nil {
		return nil, err
	}

	// Load policies
	policies, err := store.getPoliciesByEntityID(ArnUser, user.UserID)
	if err != nil {
		return nil, err
	}
	user.PolicyIDs = make([]string, len(policies))
	for i, policy := range policies {
		user.PolicyIDs[i] = policy.PolicyID
	}

	// Load groups
	groups, err := store.getGroupsByUserID(user.UserID)
	if err != nil {
		return nil, err
	}
	user.GroupIDs = make([]string, len(groups))
	for i, group := range groups {
		user.GroupIDs[i] = group.GroupID
	}

	store.Users[user.AccessKeyID.String] = &cachedUser{
		user:     &user,
		loadedAt: time.Now(),
	}

	return &user, nil
}

func (store *Store) GetUserByName(name string) (*User, error) {
	var user User

	if err := store.readDB.Get(&user, `SELECT * FROM users WHERE name = $1`, name); err != nil {
		slog.Error("get user by name", "err", err)
		return nil, err
	}

	// Load policies
	policies, err := store.getPoliciesByEntityID(ArnUser, user.UserID)
	if err != nil {
		return nil, err
	}
	user.PolicyIDs = make([]string, len(policies))
	for i, policy := range policies {
		user.PolicyIDs[i] = policy.PolicyID
	}

	// Load groups
	groups, err := store.getGroupsByUserID(user.UserID)
	if err != nil {
		return nil, err
	}
	user.GroupIDs = make([]string, len(groups))
	for i, group := range groups {
		user.GroupIDs[i] = group.GroupID
	}

	if user.AccessKeyID.Valid {
		store.Users[user.AccessKeyID.String] = &cachedUser{
			user:     &user,
			loadedAt: time.Now(),
		}
	}

	return &user, nil
}

func (store *Store) getUsersByGroupID(groupID string) ([]*User, error) {
	var users []*User

	query := `
	SELECT * FROM users WHERE user_id = (
		SELECT user_id FROM group_users WHERE group_id = $1
	)
	`

	if err := store.readDB.Select(&users, query, groupID); err != nil {
		slog.Error("get groups by user id", "err", err)
		return nil, err
	}

	return users, nil
}

func (store *Store) CreateUser(name string) (*User, error) {
	if name == config.Admin.Auth.AccessKey {
		return nil, &errs.ErrExists{Type: errs.ErrEntityTypeUser, Resource: name}
	}

	user := &User{
		UserID: GenerateUniqueID(IAMUserUniqueID),
		Name:   name,
		Arn:    toArn(ArnUser, defaultUserPath+name),
	}

	query := `
	INSERT INTO users (user_id, name, arn)
	VALUES (:user_id, :name, :arn)
	`

	_, err := store.writeDB.NamedExec(query, &user)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				slog.Warn("create user", "err", sqliteErr)
				return nil, &errs.ErrExists{Type: errs.ErrEntityTypeUser, Resource: name}
			}
		}

		slog.Error("create user", "err", err)
		return nil, err
	}

	user, err = store.GetUserByName(name)
	if err != nil {
		slog.Error("get user by name", "err", err)
		return nil, &errs.ErrNotFound{Type: errs.ErrEntityTypeUser, Resource: name}
	}

	return user, nil
}

// Creates an access key and returns the secret key in plain text
func (store *Store) CreateAccessKey(user *User) (string, error) {
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

func (store *Store) setupAdmin() error {
	encryptedData, err := security.EncryptAESGCM(config.Admin.Auth.SecretKey, config.Storage.MasterKeyPath)
	if err != nil {
		return err
	}

	store.Users[config.Admin.Auth.AccessKey] = &cachedUser{
		user: &User{
			UserID: GenerateUniqueID(IAMUserUniqueID),
			Name:   defaultAdminUserName,
			Arn:    toArn(ArnUser, defaultUserPath+defaultAdminUserName),

			AccessKeyID: sql.NullString{String: config.Admin.Auth.AccessKey, Valid: true},
			SecretKey:   encryptedData.Value,
			Salt:        encryptedData.Salt,

			GroupIDs:  []string{},
			PolicyIDs: []string{adminPolicy},
		},
		loadedAt:    time.Now(),
		neverExpire: true,
	}

	store.Policies[adminPolicy] = &cachedPolicy{
		policy: &Policy{
			PolicyID: adminPolicy,
			Document: &PolicyDocument{
				Version: latestPolicyDocumentVersion,
				Statement: []Statement{
					{
						Effect:   allow,
						Action:   []Action{Action(Any)},
						Resource: []string{Any},
					},
				},
			},
		},
		loadedAt:    time.Now(),
		neverExpire: true,
	}

	return nil
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("UserName")
	if userName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("UserName"))
		return
	}

	user, err := store.CreateUser(userName)
	if err != nil {
		var errExists *errs.ErrExists
		if errors.As(err, &errExists) {
			errs.Handle(w, errs.IAMEntityAlreadyExists(errExists.Resource))
			return
		}

		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			errs.Handle(w, errs.IAMServiceFailure())
			return
		}

		errs.Handle(w, errs.IAMServiceFailure())
		return
	}

	userPath := "/"

	response := &CreateUserResponse{
		CreateUserResult: &CreateUserResult{
			User: &UserResult{
				Path:     userPath,
				UserName: user.Name,
				UserId:   user.UserID,
				Arn:      user.Arn,
			},
			ResponseMetadata: &ResponseMetadata{
				RequestId: core.NewRequestID(),
			},
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}

func CreateAccessKeyHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("UserName")
	if userName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("UserName"))
		return
	}

	user, err := store.GetUserByName(userName)
	if err != nil {
		errs.Handle(w, errs.IAMNoSuchEntity(userName))
		return
	}

	secretKey, err := store.CreateAccessKey(user)
	if err != nil {
		slog.Error("create access key", "err", err)
		errs.Handle(w, errs.IAMServiceFailure())
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

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func ListAccessKeysHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func ListAttachedUserPoliciesHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func DeleteAccessKeyHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func DetachUserPolicyHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
