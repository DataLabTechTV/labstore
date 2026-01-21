package iam

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/IllumiKnowLabs/labstore/server/internal/core"
	"github.com/IllumiKnowLabs/labstore/server/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/server/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/server/pkg/types"
)

func (store *Store) GetUserByName(ctx context.Context, name string) (*User, error) {
	var user User

	if err := store.readDB.GetContext(ctx, &user, `SELECT * FROM users WHERE name = $1`, name); err != nil {
		slog.Error("get user by name", "err", err)
		return nil, err
	}

	// Load policies
	policies, err := store.getPoliciesByEntityID(ctx, ArnUser, user.UserID)
	if err != nil {
		return nil, err
	}
	user.PolicyIDs = make([]string, len(policies))
	for i, policy := range policies {
		user.PolicyIDs[i] = policy.PolicyID
	}

	// Load groups
	groups, err := store.getGroupsByUserID(ctx, user.UserID)
	if err != nil {
		return nil, err
	}
	user.GroupIDs = make([]string, len(groups))
	for i, group := range groups {
		user.GroupIDs[i] = group.GroupID
	}

	if user.AccessKeyID.Valid {
		store.CachedUsers[user.AccessKeyID.String] = &CachedUser{
			User:     &user,
			LoadedAt: time.Now(),
		}
	}

	return &user, nil
}

func (store *Store) GetUserByAccessKey(ctx context.Context, accessKey string) (*User, error) {
	if cachedUser, ok := store.CachedUsers[accessKey]; ok {
		if cachedUser.NeverExpire || time.Since(cachedUser.LoadedAt) < store.TTL {
			return cachedUser.User, nil
		}

		slog.Debug("invalidating cached user", "accessKey", accessKey)
		delete(store.CachedUsers, accessKey)
	}

	var user User
	if err := store.readDB.GetContext(ctx, &user, `SELECT * FROM users WHERE access_key = $1`, accessKey); err != nil {
		return nil, err
	}

	// Load policies
	policies, err := store.getPoliciesByEntityID(ctx, ArnUser, user.UserID)
	if err != nil {
		return nil, err
	}
	user.PolicyIDs = make([]string, len(policies))
	for i, policy := range policies {
		user.PolicyIDs[i] = policy.PolicyID
	}

	// Load groups
	groups, err := store.getGroupsByUserID(ctx, user.UserID)
	if err != nil {
		return nil, err
	}
	user.GroupIDs = make([]string, len(groups))
	for i, group := range groups {
		user.GroupIDs[i] = group.GroupID
	}

	store.CachedUsers[user.AccessKeyID.String] = &CachedUser{
		User:     &user,
		LoadedAt: time.Now(),
	}

	return &user, nil
}

func (store *Store) getUsersByGroupID(ctx context.Context, groupID string) ([]*User, error) {
	var users []*User

	query := `
	SELECT * FROM users WHERE user_id = (
		SELECT user_id FROM group_users WHERE group_id = $1
	)
	`

	if err := store.readDB.SelectContext(ctx, &users, query, groupID); err != nil {
		slog.Error("get groups by user id", "err", err)
		return nil, err
	}

	return users, nil
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
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

	response := &types.GetUserResponse{
		GetUserResult: &types.GetUserResult{
			User: user.Result(),
		},
		ResponseMetadata: &types.ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	helper.WriteXMLResponse(w, http.StatusOK, response)
}
