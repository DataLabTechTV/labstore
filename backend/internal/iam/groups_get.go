package iam

import (
	"log/slog"
	"net/http"
	"time"
)

func (store *Store) GetGroupByName(name string) (*Group, error) {
	var group Group

	if err := store.readDB.Get(&group, `SELECT * FROM groups WHERE name = $1`, name); err != nil {
		slog.Error("get group by name", "err", err)
		return nil, err
	}

	policies, err := store.getPoliciesByEntityID(ArnGroup, group.GroupID)
	if err != nil {
		return nil, err
	}
	group.PolicyIDs = make([]string, len(policies))
	for i, policy := range policies {
		group.PolicyIDs[i] = policy.PolicyID
	}

	store.Groups[group.Name] = &CachedGroup{
		group:    &group,
		loadedAt: time.Now(),
	}

	return &group, nil
}

func GetGroupHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
