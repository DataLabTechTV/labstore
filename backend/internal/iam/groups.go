package iam

import (
	"time"
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

type GroupResult struct {
	Path      string
	GroupName string
	GroupId   string
	Arn       string
}

func (group *Group) Result() *GroupResult {
	groupPath := "/"

	return &GroupResult{
		Path:      groupPath,
		GroupName: group.Name,
		GroupId:   group.GroupID,
		Arn:       group.Arn,
	}
}
