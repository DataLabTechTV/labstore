package iam

import (
	"database/sql"
	"encoding/xml"
	"log/slog"
	"net/http"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
	"github.com/IllumiKnowLabs/labstore/backend/internal/security"
)

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

	if user.AccessKeyID.Valid {
		if _, ok := store.CachedUsers[user.AccessKeyID.String]; ok {
			store.CachedUsers[user.AccessKeyID.String].User.AccessKeyID.String = user.AccessKeyID.String
			store.CachedUsers[user.AccessKeyID.String].User.SecretKey = user.SecretKey
			store.CachedUsers[user.AccessKeyID.String].User.Salt = user.Salt
		} else {
			store.CachedUsers[user.AccessKeyID.String] = &CachedUser{
				User:     user,
				LoadedAt: time.Now(),
			}
		}
	}

	return secretKey, nil
}

func CreateAccessKeyHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("UserName")
	if userName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("UserName"))
		return
	}

	user, err := store.GetUserByName(userName)
	if err != nil {
		errs.Handle(w, errs.IAMNoSuchEntity(string(errs.ErrEntityTypeUser), userName))
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
