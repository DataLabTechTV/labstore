package iam

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/config"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/security"
	t "github.com/IllumiKnowLabs/labstore/backend/pkg/types"
)

// Creates an access key and returns the secret key in plain text
func (store *Store) CreateAccessKey(ctx context.Context, user *User) (string, error) {
	secretKey, err := security.GeneratePassword(42)
	if err != nil {
		return "", err
	}

	user.AccessKeyID = sql.NullString{
		String: user.Name,
		Valid:  user.Name != "",
	}

	user.SecretKey, err = security.EncryptAESGCM(secretKey, config.Storage.MasterKeyPath)
	if err != nil {
		return "", err
	}

	query := `
	UPDATE
		users
	SET
		access_key = :access_key,
		secret_key = :secret_key
	WHERE
		user_id = :user_id
	`

	if _, err := store.sqlNamedExecContext(ctx, query, user); err != nil {
		return "", err
	}

	if user.AccessKeyID.Valid {
		if _, ok := store.CachedUsers[user.AccessKeyID.String]; ok {
			store.CachedUsers[user.AccessKeyID.String].User.AccessKeyID.String = user.AccessKeyID.String
			store.CachedUsers[user.AccessKeyID.String].User.SecretKey = user.SecretKey
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
	userName := r.Form.Get("UserName")
	if userName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("UserName"))
		return
	}

	ctx := r.Context()

	user, err := store.GetUserByName(ctx, userName)
	if err != nil {
		errs.Handle(w, errs.IAMNoSuchEntity(string(errs.ErrEntityTypeUser), userName))
		return
	}

	secretKey, err := store.CreateAccessKey(ctx, user)
	if err != nil {
		slog.Error("create access key", "err", err)
		errs.Handle(w, errs.IAMServiceFailure())
		return
	}

	response := &t.CreateAccessKeyResponse{
		CreateAccessKeyResult: &t.CreateAccessKeyResult{
			AccessKey: &t.AccessKeyResult{
				UserName:        user.Name,
				AccessKeyId:     user.AccessKeyID.String,
				Status:          t.AccessKeyActive,
				SecretAccessKey: secretKey,
			},
			ResponseMetadata: &t.ResponseMetadata{
				RequestId: core.NewRequestID(),
			},
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}
