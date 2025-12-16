package iam

import (
	"encoding/xml"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

const (
	defaultGroupPath = "/"
)

type Group struct {
	GroupID string `db:"group_id"`
	Name    string `db:"name"`
	Arn     string `db:"arn"`

	UserIDs   []string
	PolicyIDs []string
}

type cachedGroup struct {
	group       *Group
	loadedAt    time.Time
	neverExpire bool
}

type CreateGroupResponse struct {
	XMLName           xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ CreateGroupResponse"`
	CreateGroupResult *CreateGroupResult
}

type CreateGroupResult struct {
	Group            *GroupResult
	ResponseMetadata *ResponseMetadata
}

type GroupResult struct {
	Path      string
	GroupName string
	GroupId   string
	Arn       string
}

type AddUserToGroupResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ AddUserToGroupResponse"`
	ResponseMetadata *ResponseMetadata
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

	store.Groups[groupID] = &cachedGroup{
		group:    &group,
		loadedAt: time.Now(),
	}

	return &group, nil
}

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

	store.Groups[group.Name] = &cachedGroup{
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

func (store *Store) CreateGroup(name string) (*Group, error) {
	group := &Group{
		GroupID: GenerateUniqueID(IAMGroupUniqueID),
		Name:    name,
		Arn:     toArn(ArnGroup, defaultGroupPath+name),
	}

	query := `
	INSERT INTO groups (group_id, name, arn)
	VALUES (:group_id, :name, :arn)
	`

	_, err := store.writeDB.NamedExec(query, &group)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				slog.Warn("create group", "err", sqliteErr)
				return nil, &errs.ErrExists{Type: errs.ErrEntityTypeGroup, Resource: name}
			}
		}

		slog.Error("create group", "err", err)
		return nil, err
	}

	group, err = store.GetGroupByName(name)
	if err != nil {
		slog.Error("get group by name", "err", err)
		return nil, &errs.ErrNotFound{Type: errs.ErrEntityTypeGroup, Resource: name}
	}

	return group, nil
}

func (store *Store) AddUserToGroup(userName, groupName string) error {
	user, err := store.GetUserByName(userName)
	if err != nil {
		slog.Error("get user by name", "err", err)
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeUser, Resource: userName}
	}

	group, err := store.GetGroupByName(groupName)
	if err != nil {
		slog.Error("get group by name", "err", err)
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeGroup, Resource: groupName}
	}

	query := `
		INSERT INTO group_users (group_id, user_id)
		VALUES ($1, $2)
	`
	_, err = store.writeDB.Exec(query, group.GroupID, user.UserID)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY {
				slog.Warn("add user to group", "userName", userName, "groupName", groupName, "err", sqliteErr)
				return nil
			}
		}

		slog.Error("add user to group", "userName", userName, "groupName", groupName, "err", sqliteErr)
		return err
	}

	user, err = store.GetUserByName(userName)
	if err != nil {
		return err
	}

	group.UserIDs = append(group.UserIDs, user.UserID)

	return nil
}

func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	groupName := r.URL.Query().Get("GroupName")
	if groupName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("GroupName"))
		return
	}

	group, err := store.CreateGroup(groupName)
	if err != nil {
		var errExists *errs.ErrExists
		if errors.As(err, &errExists) {
			errs.Handle(w, errs.IAMEntityAlreadyExists(errExists.Resource))
			return
		}

		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			errs.Handle(w, errs.IAMServiceFailure())
			return
		}

		errs.Handle(w, errs.IAMServiceFailure())
		return
	}

	groupPath := "/"

	response := &CreateGroupResponse{
		CreateGroupResult: &CreateGroupResult{
			Group: &GroupResult{
				Path:      groupPath,
				GroupName: group.Name,
				GroupId:   group.GroupID,
				Arn:       group.Arn,
			},
			ResponseMetadata: &ResponseMetadata{
				RequestId: core.NewRequestID(),
			},
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}

func AddUserToGroupHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("UserName")
	if userName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("UserName"))
		return
	}

	groupName := r.URL.Query().Get("GroupName")
	if groupName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("GroupName"))
		return
	}

	if err := store.AddUserToGroup(userName, groupName); err != nil {
		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			errs.Handle(w, errs.IAMNoSuchEntity(errNotFound.Resource))
			return
		}

		errs.Handle(w, errs.IAMServiceFailure())
		return
	}

	response := &AddUserToGroupResponse{
		ResponseMetadata: &ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}

func GetGroupHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func ListAttachedGroupPoliciesHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func DeleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func RemoveUserFromGroupHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func DetachGroupPolicyHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
