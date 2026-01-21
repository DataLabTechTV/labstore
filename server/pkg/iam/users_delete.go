package iam

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/server/internal/core"
	"github.com/IllumiKnowLabs/labstore/server/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/server/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/server/pkg/types"
)

func (store *Store) DeleteUser(ctx context.Context, name string) error {
	if name == defaultAdminUserName {
		return &errs.ErrForbidden{Type: errs.ErrEntityTypeUser, Resource: name}
	}

	user, err := store.GetUserByName(ctx, name)
	if err != nil {
		slog.Error("delete user", "err", err)
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeUser, Resource: name}
	}

	query := `
	DELETE FROM users
	WHERE name = $1
	`

	_, err = store.sqlExecContext(ctx, query, user.Name)
	if err != nil {
		slog.Error("delete user", "err", err)
		return err
	}

	delete(store.CachedUsers, user.AccessKeyID.String)

	return nil
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.Form.Get("UserName")
	if userName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("UserName"))
		return
	}

	ctx := r.Context()

	if err := store.DeleteUser(ctx, userName); err != nil {
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

	response := &types.DeleteUserResponse{
		ResponseMetadata: &types.ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	helper.WriteXMLResponse(w, http.StatusOK, response)
}
