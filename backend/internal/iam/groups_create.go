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

func (store *Store) CreateGroup(ctx context.Context, name string) (*Group, error) {
	group := &Group{
		GroupID: GenerateUniqueID(IAMGroupUniqueID),
		Name:    name,
		Arn:     toArn(ArnGroup, defaultGroupPath+name),
	}

	query := `
	INSERT INTO groups (group_id, name, arn)
	VALUES (:group_id, :name, :arn)
	`

	_, err := store.sqlNamedExecContext(ctx, query, &group)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				slog.Warn("create group", "err", sqliteErr)
				return nil, &errs.ErrExists{Type: errs.ErrEntityTypeGroup, Resource: name}
			}
		}

		slog.Error("create group", "err", err)
		return nil, err
	}

	group, err = store.GetGroupByName(ctx, name)
	if err != nil {
		slog.Error("get group by name", "err", err)
		return nil, &errs.ErrNotFound{Type: errs.ErrEntityTypeGroup, Resource: name}
	}

	return group, nil
}

func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	groupName := r.Form.Get("GroupName")
	if groupName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("GroupName"))
		return
	}

	ctx := r.Context()

	group, err := store.CreateGroup(ctx, groupName)
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

	response := &types.CreateGroupResponse{
		CreateGroupResult: &types.CreateGroupResult{
			Group: group.Result(),
			ResponseMetadata: &types.ResponseMetadata{
				RequestId: core.NewRequestID(),
			},
		},
	}

	helper.WriteXMLResponse(w, http.StatusOK, response)
}
