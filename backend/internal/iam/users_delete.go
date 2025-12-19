package iam

import (
	"context"
	"encoding/xml"
	"errors"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
)

type DeleteUserResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ DeleteUserResponse"`
	ResponseMetadata *ResponseMetadata
}

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

	response := &DeleteUserResponse{
		ResponseMetadata: &ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}
