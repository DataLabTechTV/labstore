package iam

import (
	"context"
	"encoding/xml"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
)

type DetachUserPolicyResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ DetachUserPolicyResponse"`
	ResponseMetadata *ResponseMetadata
}

type DetachGroupPolicyResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ DetachGroupPolicyResponse"`
	ResponseMetadata *ResponseMetadata
}

func (store *Store) DetachUserPolicy(ctx context.Context, userName, policyArn string) error {
	if userName == defaultAdminUserName {
		return &errs.ErrForbidden{Type: errs.ErrEntityTypeUser, Resource: userName}
	}

	user, err := store.GetUserByName(ctx, userName)
	if err != nil {
		slog.Error("detach user policy", "err", err)
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeAccessKey, Resource: userName}
	}

	policy, err := store.GetPolicyByArn(ctx, policyArn)
	if err != nil {
		slog.Error("detach user policy", "err", err)
		return &errs.ErrNotFound{Type: errs.ErrEntityTypePolicy, Resource: policyArn}
	}

	query := `
	DELETE FROM user_policies
	WHERE user_id = $1
	AND policy_id = $2
	`

	res, err := store.sqlExecContext(ctx, query, user.UserID, policy.PolicyID)
	if err != nil {
		slog.Error("detach user policy", "err", err)
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		slog.Error("detach user policy", "err", err)
		return err
	}
	if n < 1 {
		slog.Error("detach user policy: 0 deleted")
		return &errs.ErrUserPolicyNotAttached{UserID: user.UserID, PolicyID: policy.PolicyID}
	}

	user, err = store.GetUserByName(ctx, userName)
	if err != nil {
		slog.Error("detach user policy", "err", err)
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeUser, Resource: userName}
	}

	if user.AccessKeyID.Valid {
		if _, ok := store.CachedUsers[user.AccessKeyID.String]; ok {
			store.CachedUsers[user.AccessKeyID.String].User = user
			store.CachedUsers[user.AccessKeyID.String].LoadedAt = time.Now()
		}
	}

	return nil
}

func (store *Store) DetachGroupPolicy(ctx context.Context, groupName, policyArn string) error {
	group, err := store.GetGroupByName(ctx, groupName)
	if err != nil {
		slog.Error("detach group policy", "err", err)
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeGroup, Resource: groupName}
	}

	policy, err := store.GetPolicyByArn(ctx, policyArn)
	if err != nil {
		slog.Error("detach group policy", "err", err)
		return &errs.ErrNotFound{Type: errs.ErrEntityTypePolicy, Resource: policyArn}
	}

	query := `
	DELETE FROM group_policies
	WHERE group_id = $1
	AND policy_id = $2
	`

	res, err := store.sqlExecContext(ctx, query, group.GroupID, policy.PolicyID)
	if err != nil {
		slog.Error("detach group policy", "err", err)
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		slog.Error("detach group policy", "err", err)
		return err
	}
	if n < 1 {
		slog.Error("detach group policy: 0 deleted")
		return &errs.ErrGroupPolicyNotAttached{GroupID: group.GroupID, PolicyID: policy.PolicyID}
	}

	group, err = store.GetGroupByName(ctx, groupName)
	if err != nil {
		slog.Error("detach group policy", "err", err)
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeGroup, Resource: groupName}
	}

	if _, ok := store.CachedGroups[group.GroupID]; ok {
		store.CachedGroups[group.GroupID].Group = group
		store.CachedGroups[group.GroupID].LoadedAt = time.Now()
	}

	return nil
}

func DetachUserPolicyHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("UserName")
	if userName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("UserName"))
		return
	}

	policyArn := r.URL.Query().Get("PolicyArn")
	if userName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("PolicyArn"))
		return
	}

	ctx := r.Context()

	if err := store.DetachUserPolicy(ctx, userName, policyArn); err != nil {
		var errForbidden *errs.ErrForbidden
		if errors.As(err, &errForbidden) {
			errs.Handle(w, errs.IAMDeleteConflict(errForbidden.Resource))
			return
		}

		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			errs.Handle(w, errs.IAMNoSuchEntity(string(errNotFound.Type), errNotFound.Resource))
			return
		}

		var errUserPolicyNotAttached *errs.ErrUserPolicyNotAttached
		if errors.As(err, &errUserPolicyNotAttached) {
			errs.Handle(w, errs.IAMNoSuchEntity(
				errs.ErrEntityTypeUserPolicy,
				errUserPolicyNotAttached.Description(),
			))
			return
		}

		errs.Handle(w, errs.IAMServiceFailure())
		return
	}

	response := &DetachUserPolicyResponse{
		ResponseMetadata: &ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}

func DetachGroupPolicyHandler(w http.ResponseWriter, r *http.Request) {
	groupName := r.URL.Query().Get("GroupName")
	if groupName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("GroupName"))
		return
	}

	policyArn := r.URL.Query().Get("PolicyArn")
	if groupName == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("PolicyArn"))
		return
	}

	ctx := r.Context()

	if err := store.DetachGroupPolicy(ctx, groupName, policyArn); err != nil {
		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			errs.Handle(w, errs.IAMNoSuchEntity(string(errNotFound.Type), errNotFound.Resource))
			return
		}

		var errGroupPolicyNotAttached *errs.ErrGroupPolicyNotAttached
		if errors.As(err, &errGroupPolicyNotAttached) {
			errs.Handle(w, errs.IAMNoSuchEntity(
				errs.ErrEntityTypeGroupPolicy,
				errGroupPolicyNotAttached.Description(),
			))
			return
		}

		errs.Handle(w, errs.IAMServiceFailure())
		return
	}

	response := &DetachUserPolicyResponse{
		ResponseMetadata: &ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}
