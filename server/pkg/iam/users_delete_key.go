package iam

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/server/internal/core"
	"github.com/IllumiKnowLabs/labstore/server/pkg/config"
	"github.com/IllumiKnowLabs/labstore/server/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/server/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/server/pkg/types"
)

func (store *Store) DeleteAccessKey(ctx context.Context, userName, accessKeyID string) error {
	if userName == defaultAdminUserName {
		return &errs.ErrForbidden{Type: errs.ErrEntityTypeUser, Resource: userName}
	}

	if accessKeyID == config.Admin.Auth.AccessKey {
		return &errs.ErrForbidden{Type: errs.ErrEntityTypeAccessKey, Resource: accessKeyID}
	}

	user, err := store.GetUserByAccessKey(ctx, accessKeyID)
	if err != nil {
		slog.Error("delete access key", "err", err)
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeAccessKey, Resource: accessKeyID}
	}

	query := `
	UPDATE
		users
	SET
		access_key = NULL,
		secret_key = NULL
	WHERE
		name = :name
		AND access_key = :access_key
	`

	_, err = store.sqlNamedExecContext(ctx, query, &user)
	if err != nil {
		slog.Error("delete user", "err", err)
		return err
	}

	if user.AccessKeyID.Valid {
		delete(store.CachedUsers, user.AccessKeyID.String)
	}

	return nil
}

func DeleteAccessKeyHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.Form.Get("UserName")
	if userName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("UserName"))
		return
	}

	accessKeyId := r.Form.Get("AccessKeyId")
	if userName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("AccessKeyId"))
		return
	}

	ctx := r.Context()

	if err := store.DeleteAccessKey(ctx, userName, accessKeyId); err != nil {
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

	response := &types.DeleteAccessKeyResponse{
		ResponseMetadata: &types.ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	helper.WriteXMLResponse(w, http.StatusOK, response)
}
