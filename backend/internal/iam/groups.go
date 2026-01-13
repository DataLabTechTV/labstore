package iam

import (
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/types"
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

func (group *Group) Result() *types.GroupResult {
	groupPath := "/"

	return &types.GroupResult{
		Path:      groupPath,
		GroupName: group.Name,
		GroupId:   group.GroupID,
		Arn:       group.Arn,
	}
}
