package iam

import (
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

func (store *Store) DetachUserPolicy(userName, policyArn string) error {
	if userName == defaultAdminUserName {
		return &errs.ErrForbidden{Type: errs.ErrEntityTypeUser, Resource: userName}
	}

	user, err := store.GetUserByName(userName)
	if err != nil {
		slog.Error("detach user policy", "err", err)
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeAccessKey, Resource: userName}
	}

	policy, err := store.GetPolicyByArn(policyArn)
	if err != nil {
		slog.Error("detach user policy", "err", err)
		return &errs.ErrNotFound{Type: errs.ErrEntityTypePolicy, Resource: policyArn}
	}

	query := `
	DELETE FROM user_policies
	WHERE user_id = $1
	AND policy_id = $2
	`

	res, err := store.writeDB.Exec(query, user.UserID, policy.PolicyID)
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
		return &errs.ErrNotFound{Type: errs.ErrEntityTypeUserPolicy, Resource: policyArn}
	}

	user, err = store.GetUserByName(userName)
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

	if err := store.DetachUserPolicy(userName, policyArn); err != nil {
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
			errs.Handle(w, errs.IAMNoSuchEntity(errs.ErrEntityTypeUserPolicy, policyArn))
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
