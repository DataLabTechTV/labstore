package iam

import (
	"time"

	t "github.com/IllumiKnowLabs/labstore/backend/pkg/types"
)

const (
	defaultGroupPath = "/"
)

type CachedGroup struct {
	Group       *Group
	LoadedAt    time.Time
	NeverExpire bool
}

type Group struct {
	GroupID string `db:"group_id"`
	Name    string `db:"name"`
	Arn     string `db:"arn"`

	UserIDs   []string
	PolicyIDs []string
}

func (group *Group) Result() *t.GroupResult {
	groupPath := "/"

	return &t.GroupResult{
		Path:      groupPath,
		GroupName: group.Name,
		GroupId:   group.GroupID,
		Arn:       group.Arn,
	}
}
