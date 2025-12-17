package iam

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type AttachUserPolicyResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ AttachUserPolicyResponse"`
	ResponseMetadata *ResponseMetadata
}

func (store *Store) AttachPolicy(arnType ArnType, policyArn, entityName string) error {
	var entity any

	policy, err := store.GetPolicyByArn(policyArn)
	if err != nil {
		slog.Error("get policy by arn", "err", err)
		return &errs.ErrNotFound{Type: errs.ErrEntityTypePolicy, Resource: policyArn}
	}

	var tableName string
	var idFieldName string
	var idFieldValue string

	switch arnType {
	case ArnUser:
		user, err := store.GetUserByName(entityName)
		if err != nil {
			slog.Error("get user by name", "err", err)
			return &errs.ErrNotFound{Type: errs.ErrEntityTypeUser, Resource: entityName}
		}

		tableName = "user_policies"
		idFieldName = "user_id"
		idFieldValue = user.UserID
		entity = user
	case ArnGroup:
		group, err := store.GetGroupByName(entityName)
		if err != nil {
			slog.Error("get group by name", "err", err)
			return &errs.ErrNotFound{Type: errs.ErrEntityTypeGroup, Resource: entityName}
		}

		tableName = "group_policies"
		idFieldName = "group_id"
		idFieldValue = group.GroupID
		entity = group
	default:
		return errors.New("unsupported arn type")
	}

	query_tmpl := `
		INSERT INTO %s (%s, policy_id)
		VALUES ($1, $2)
	`
	query := fmt.Sprintf(query_tmpl, tableName, idFieldName)

	_, err = store.sqlExec(query, idFieldValue, policy.PolicyID)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY {
				slog.Warn("attach policy", "arnType", arnType, "err", sqliteErr)
				return nil
			}
		}

		slog.Error("attach policy", "arnType", arnType, "err", err)
		return err
	}

	switch e := entity.(type) {
	case *User:
		if e.AccessKeyID.Valid {
			if _, ok := store.CachedUsers[e.AccessKeyID.String]; ok {
				store.CachedUsers[e.AccessKeyID.String].User.PolicyIDs = append(e.PolicyIDs, policy.PolicyID)
			}
		}
	case *Group:
		if _, ok := store.CachedGroups[e.GroupID]; ok {
			store.CachedGroups[e.GroupID].Group.PolicyIDs = append(e.PolicyIDs, policy.PolicyID)
		}
	}

	return nil
}

func AttachUserPolicyHandler(w http.ResponseWriter, r *http.Request) {
	policyArn := r.URL.Query().Get("PolicyArn")
	if policyArn == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("PolicyArn"))
		return
	}

	userName := r.URL.Query().Get("UserName")
	if userName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("UserName"))
		return
	}

	if err := store.AttachPolicy(ArnUser, policyArn, userName); err != nil {
		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			errs.Handle(w, errs.IAMNoSuchEntity(string(errNotFound.Type), errNotFound.Resource))
			return
		}

		errs.Handle(w, errs.IAMServiceFailure())
		return
	}

	response := &AttachUserPolicyResponse{
		ResponseMetadata: &ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}

func AttachGroupPolicyHandler(w http.ResponseWriter, r *http.Request) {
	policyArn := r.URL.Query().Get("PolicyArn")
	if policyArn == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("PolicyArn"))
		return
	}

	groupName := r.URL.Query().Get("GroupName")
	if groupName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("GroupName"))
		return
	}

	if err := store.AttachPolicy(ArnGroup, policyArn, groupName); err != nil {
		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			errs.Handle(w, errs.IAMNoSuchEntity(string(errNotFound.Type), errNotFound.Resource))
			return
		}

		errs.Handle(w, errs.IAMServiceFailure())
		return
	}

	response := &AttachUserPolicyResponse{
		ResponseMetadata: &ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}
