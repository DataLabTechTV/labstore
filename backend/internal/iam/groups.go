package iam

import (
	"log/slog"
	"time"
)

const (
	defaultGroupPath = "/"
)

type CachedGroup struct {
	group       *Group
	loadedAt    time.Time
	neverExpire bool
}

type Group struct {
	GroupID string `db:"group_id"`
	Name    string `db:"name"`
	Arn     string `db:"arn"`

	UserIDs   []string
	PolicyIDs []string
}

func (store *Store) GetGroupByID(groupID string) (*Group, error) {
	if cachedGroup, ok := store.Groups[groupID]; ok {
		if cachedGroup.neverExpire || time.Since(cachedGroup.loadedAt) < store.ttl {
			return cachedGroup.group, nil
		}

		slog.Debug("invalidating cached group", "groupID", groupID)
		delete(store.Groups, groupID)
	}

	var group Group

	query := `SELECT * FROM groups WHERE group_id = $1`
	if err := store.readDB.Get(&group, query, groupID); err != nil {
		return nil, err
	}

	// Load policies
	policies, err := store.getPoliciesByEntityID(ArnGroup, group.GroupID)
	if err != nil {
		return nil, err
	}
	group.PolicyIDs = make([]string, len(policies))
	for i, policy := range policies {
		group.PolicyIDs[i] = policy.PolicyID
	}

	// Load groups
	users, err := store.getUsersByGroupID(group.GroupID)
	if err != nil {
		return nil, err
	}
	group.UserIDs = make([]string, len(users))
	for i, user := range users {
		group.UserIDs[i] = user.UserID
	}

	store.Groups[groupID] = &CachedGroup{
		group:    &group,
		loadedAt: time.Now(),
	}

	return &group, nil
}

func (store *Store) getGroupsByUserID(userID string) ([]*Group, error) {
	var groups []*Group

	query := `
	SELECT * FROM groups WHERE group_id = (
		SELECT group_id FROM group_users WHERE user_id = $1
	)
	`

	if err := store.readDB.Select(&groups, query, userID); err != nil {
		slog.Error("get groups by user id", "err", err)
		return nil, err
	}

	return groups, nil
}
