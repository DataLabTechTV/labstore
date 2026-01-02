package iam

import (
	"context"
	"encoding/xml"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
	t "github.com/IllumiKnowLabs/labstore/backend/pkg/types"
)

type RemoveUserFromGroupResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ RemoveUserFromGroupResponse"`
	ResponseMetadata *t.ResponseMetadata
}

func (store *Store) RemoveUserFromGroup(ctx context.Context, userName, groupName string) error {
	group, err := store.GetGroupByName(ctx, groupName)
	if err != nil {
		slog.Error("remove user from group", "err", err)
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeGroup, Resource: groupName}
	}

	user, err := store.GetUserByName(ctx, userName)
	if err != nil {
		slog.Error("remove user from group", "err", err)
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeUser, Resource: userName}
	}

	query := `
	DELETE FROM group_users
	WHERE group_id = $1
	AND user_id = $2
	`

	res, err := store.sqlExecContext(ctx, query, group.GroupID, user.UserID)
	if err != nil {
		slog.Error("remove user from group", "err", err)
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		slog.Error("remove user from group", "err", err)
		return err
	}
	if n < 1 {
		slog.Error("remove user from group: 0 deleted")
		return &errs.ErrUserNotInGroup{UserID: user.UserID, GroupID: group.GroupID}
	}

	group, err = store.GetGroupByName(ctx, groupName)
	if err != nil {
		slog.Error("remove user from group", "err", err)
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeGroup, Resource: groupName}
	}

	if _, ok := store.CachedGroups[group.GroupID]; ok {
		store.CachedGroups[group.GroupID].Group = group
		store.CachedGroups[group.GroupID].LoadedAt = time.Now()
	}

	return nil
}

func RemoveUserFromGroupHandler(w http.ResponseWriter, r *http.Request) {
	groupName := r.Form.Get("GroupName")
	if groupName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("GroupName"))
		return
	}

	userName := r.Form.Get("UserName")
	if groupName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("UserName"))
		return
	}

	ctx := r.Context()

	if err := store.RemoveUserFromGroup(ctx, userName, groupName); err != nil {
		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			errs.Handle(w, errs.IAMNoSuchEntity(string(errNotFound.Type), errNotFound.Resource))
			return
		}

		var errUserNotInGroup *errs.ErrUserNotInGroup
		if errors.As(err, &errUserNotInGroup) {
			errs.Handle(w, errs.IAMNoSuchEntity(
				errs.ErrEntityTypeUserGroup,
				errUserNotInGroup.Description(),
			))
			return
		}

		errs.Handle(w, errs.IAMServiceFailure())
		return
	}

	response := &RemoveUserFromGroupResponse{
		ResponseMetadata: &t.ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}
