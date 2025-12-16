package iam

import "time"

type Group struct {
	GroupID string `db:"group_id"`
	Name    string `db:"name"`

	UserIDs   []string
	PolicyIDs []string
}

//nolint:unused
type cachedGroup struct {
	group       *Group
	loadedAt    time.Time
	neverExpire bool
}
