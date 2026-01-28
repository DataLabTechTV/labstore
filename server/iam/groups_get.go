package iam

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/IllumiKnowLabs/labstore/server/core"
	"github.com/IllumiKnowLabs/labstore/server/errs"
	"github.com/IllumiKnowLabs/labstore/server/helper"
	"github.com/IllumiKnowLabs/labstore/server/types"
)

func (store *Store) GetGroupByName(ctx context.Context, name string) (*Group, error) {
	var group Group

	if err := store.readDB.GetContext(ctx, &group, `SELECT * FROM groups WHERE name = $1`, name); err != nil {
		slog.Error("get group by name", "err", err)
		return nil, err
	}

	// Load policies
	policies, err := store.getPoliciesByEntityID(ctx, ArnGroup, group.GroupID)
	if err != nil {
		return nil, err
	}
	group.PolicyIDs = make([]string, len(policies))
	for i, policy := range policies {
		group.PolicyIDs[i] = policy.PolicyID
	}

	// Load users
	users, err := store.getUsersByGroupID(ctx, group.GroupID)
	if err != nil {
		return nil, err
	}
	group.UserIDs = make([]string, len(users))
	for i, user := range users {
		group.UserIDs[i] = user.UserID
	}

	store.CachedGroups[group.Name] = &CachedGroup{
		Group:    &group,
		LoadedAt: time.Now(),
	}

	return &group, nil
}

func (store *Store) GetGroupByID(ctx context.Context, groupID string) (*Group, error) {
	if cachedGroup, ok := store.CachedGroups[groupID]; ok {
		if cachedGroup.NeverExpire || time.Since(cachedGroup.LoadedAt) < store.TTL {
			return cachedGroup.Group, nil
		}

		slog.Debug("invalidating cached group", "groupID", groupID)
		delete(store.CachedGroups, groupID)
	}

	var group Group

	query := `SELECT * FROM groups WHERE group_id = $1`
	if err := store.readDB.GetContext(ctx, &group, query, groupID); err != nil {
		return nil, err
	}

	// Load policies
	policies, err := store.getPoliciesByEntityID(ctx, ArnGroup, group.GroupID)
	if err != nil {
		return nil, err
	}
	group.PolicyIDs = make([]string, len(policies))
	for i, policy := range policies {
		group.PolicyIDs[i] = policy.PolicyID
	}

	// Load groups
	users, err := store.getUsersByGroupID(ctx, group.GroupID)
	if err != nil {
		return nil, err
	}
	group.UserIDs = make([]string, len(users))
	for i, user := range users {
		group.UserIDs[i] = user.UserID
	}

	store.CachedGroups[groupID] = &CachedGroup{
		Group:    &group,
		LoadedAt: time.Now(),
	}

	return &group, nil
}

func (store *Store) getGroupsByUserID(ctx context.Context, userID string) ([]*Group, error) {
	var groups []*Group

	query := `
	SELECT * FROM groups WHERE group_id = (
		SELECT group_id FROM group_users WHERE user_id = $1
	)
	`

	if err := store.readDB.SelectContext(ctx, &groups, query, userID); err != nil {
		slog.Error("get groups by user id", "err", err)
		return nil, err
	}

	return groups, nil
}

func GetGroupHandler(w http.ResponseWriter, r *http.Request) {
	groupName := r.Form.Get("GroupName")
	if groupName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("GroupName"))
		return
	}

	ctx := r.Context()

	group, err := store.GetGroupByName(ctx, groupName)
	if err != nil {
		errs.Handle(w, errs.IAMNoSuchEntity(string(errs.ErrEntityTypeGroup), groupName))
		return
	}

	users, err := store.getUsersByGroupID(ctx, group.GroupID)
	if err != nil {
		errs.Handle(w, errs.IAMServiceFailure())
		return
	}

	userResults := make([]*types.UserResult, len(users))
	for i, user := range users {
		userResults[i] = user.Result()
	}

	response := &types.GetGroupResponse{
		GetGroupResult: &types.GetGroupResult{
			Group: group.Result(),
			Users: &types.UserMembers{
				Member: userResults,
			},
			IsTruncated: false,
		},
		ResponseMetadata: &types.ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	helper.WriteXMLResponse(w, http.StatusOK, response)
}
