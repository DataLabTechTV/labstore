package iam

import (
	"encoding/xml"
	"log/slog"
	"net/http"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
)

type GetGroupResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ GetGroupResponse"`
	GetGroupResult   *GetGroupResult
	ResponseMetadata *ResponseMetadata
}

type GetGroupResult struct {
	Group       *GroupResult
	Users       *UserMembers
	IsTruncated bool
}

type UserMembers struct {
	Member []*UserResult `xml:"member"`
}

func (store *Store) GetGroupByName(name string) (*Group, error) {
	var group Group

	if err := store.readDB.Get(&group, `SELECT * FROM groups WHERE name = $1`, name); err != nil {
		slog.Error("get group by name", "err", err)
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

	// Load users
	users, err := store.getUsersByGroupID(group.GroupID)
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

func GetGroupHandler(w http.ResponseWriter, r *http.Request) {
	groupName := r.URL.Query().Get("GroupName")
	if groupName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("GroupName"))
		return
	}

	group, err := store.GetGroupByName(groupName)
	if err != nil {
		errs.Handle(w, errs.IAMNoSuchEntity(string(errs.ErrEntityTypeGroup), groupName))
		return
	}

	users, err := store.getUsersByGroupID(group.GroupID)
	if err != nil {
		errs.Handle(w, errs.IAMServiceFailure())
		return
	}

	userResults := make([]*UserResult, len(users))
	for i, user := range users {
		userResults[i] = user.Result()
	}

	response := &GetGroupResponse{
		GetGroupResult: &GetGroupResult{
			Group: group.Result(),
			Users: &UserMembers{
				Member: userResults,
			},
			IsTruncated: false,
		},
		ResponseMetadata: &ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}
