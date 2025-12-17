package iam

import (
	"database/sql"
	"encoding/xml"
	"errors"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
)

type DeleteAccessKeyResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ DeleteAccessKeyResponse"`
	ResponseMetadata *ResponseMetadata
}

func (store *Store) DeleteAccessKey(userName, accessKeyID string) error {
	if userName == defaultAdminUserName {
		return &errs.ErrForbidden{Type: errs.ErrEntityTypeUser, Resource: userName}
	}

	if accessKeyID == config.Admin.Auth.AccessKey {
		return &errs.ErrForbidden{Type: errs.ErrEntityTypeAccessKey, Resource: accessKeyID}
	}

	user, err := store.GetUserByAccessKey(accessKeyID)
	if err != nil {
		slog.Error("delete access key", "err", err)
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeAccessKey, Resource: accessKeyID}
	}

	query := `
	UPDATE users
	SET
		access_key = NULL,
		secret_key = NULL,
		salt = NUll
	WHERE name = :name AND access_key = :access_key
	`

	_, err = store.writeDB.NamedExec(query, &user)
	if err != nil {
		slog.Error("delete user", "err", err)
		return err
	}

	if user.AccessKeyID.Valid {
		if _, ok := store.CachedUsers[user.AccessKeyID.String]; ok {
			store.CachedUsers[user.AccessKeyID.String].User.AccessKeyID = sql.NullString{Valid: false}
			store.CachedUsers[user.AccessKeyID.String].User.SecretKey = nil
			store.CachedUsers[user.AccessKeyID.String].User.Salt = nil
		}
	}

	return nil
}

func DeleteAccessKeyHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("UserName")
	if userName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("UserName"))
		return
	}

	accessKeyId := r.URL.Query().Get("AccessKeyId")
	if userName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("AccessKeyId"))
		return
	}

	if err := store.DeleteAccessKey(userName, accessKeyId); err != nil {
		var errForbidden *errs.ErrForbidden
		if errors.As(err, &errForbidden) {
			errs.Handle(w, errs.IAMDeleteConflict(errForbidden.Resource))
			return
		}

		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			errs.Handle(w, errs.IAMNoSuchEntity(string(errNotFound.Type), errNotFound.Resource))
			return
		}

		errs.Handle(w, errs.IAMServiceFailure())
		return
	}

	response := &DeleteAccessKeyResponse{
		ResponseMetadata: &ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}
