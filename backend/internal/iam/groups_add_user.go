package iam

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/backend/pkg/types"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func (store *Store) AddUserToGroup(ctx context.Context, userName, groupName string) error {
	user, err := store.GetUserByName(ctx, userName)
	if err != nil {
		slog.Error("get user by name", "err", err)
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeUser, Resource: userName}
	}

	group, err := store.GetGroupByName(ctx, groupName)
	if err != nil {
		slog.Error("get group by name", "err", err)
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeGroup, Resource: groupName}
	}

	query := `
		INSERT INTO group_users (group_id, user_id)
		VALUES ($1, $2)
	`
	_, err = store.sqlExecContext(ctx, query, group.GroupID, user.UserID)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY {
				slog.Warn("add user to group", "userName", userName, "groupName", groupName, "err", sqliteErr)
				return nil
			}
		}

		slog.Error("add user to group", "userName", userName, "groupName", groupName, "err", sqliteErr)
		return err
	}

	user, err = store.GetUserByName(ctx, userName)
	if err != nil {
		return err
	}

	group.UserIDs = append(group.UserIDs, user.UserID)

	return nil
}

func AddUserToGroupHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.Form.Get("UserName")
	if userName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("UserName"))
		return
	}

	groupName := r.Form.Get("GroupName")
	if groupName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("GroupName"))
		return
	}

	ctx := r.Context()

	if err := store.AddUserToGroup(ctx, userName, groupName); err != nil {
		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			errs.Handle(w, errs.IAMNoSuchEntity(string(errNotFound.Type), errNotFound.Resource))
			return
		}

		errs.Handle(w, errs.IAMServiceFailure())
		return
	}

	response := &types.AddUserToGroupResponse{
		ResponseMetadata: &types.ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	helper.WriteXMLResponse(w, http.StatusOK, response)
}
