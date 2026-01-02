package iam

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	t "github.com/IllumiKnowLabs/labstore/backend/pkg/types"
)

func (store *Store) DeleteGroup(ctx context.Context, name string) error {
	group, err := store.GetGroupByName(ctx, name)
	if err != nil {
		slog.Error("delete group", "err", err)
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeGroup, Resource: name}
	}

	query := `
	DELETE FROM groups
	WHERE name = $1
	`

	_, err = store.sqlExecContext(ctx, query, group.Name)
	if err != nil {
		slog.Error("delete group", "err", err)
		return err
	}

	delete(store.CachedGroups, group.GroupID)

	return nil
}

func DeleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	groupName := r.Form.Get("GroupName")
	if groupName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("GroupName"))
		return
	}

	ctx := r.Context()

	if err := store.DeleteGroup(ctx, groupName); err != nil {
		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			errs.Handle(w, errs.IAMNoSuchEntity(string(errNotFound.Type), errNotFound.Resource))
			return
		}

		errs.Handle(w, errs.IAMServiceFailure())
		return
	}

	response := &t.DeleteGroupResponse{
		ResponseMetadata: &t.ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	helper.WriteXMLResponse(w, http.StatusOK, response)
}
