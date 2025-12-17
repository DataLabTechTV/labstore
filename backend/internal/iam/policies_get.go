package iam

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
)

type GetPolicyResponse struct {
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ GetPolicyResponse"`
	GetPolicyResult  *GetPolicyResult
	ResponseMetadata *ResponseMetadata
}

type GetPolicyResult struct {
	Policy *PolicyResult
}

func (store *Store) GetPolicyByArn(arn string) (*Policy, error) {
	var policy Policy

	query := `SELECT * FROM policies WHERE arn = $1`
	if err := store.readDB.Get(&policy, query, arn); err != nil {
		return nil, err
	}

	attachments, err := store.countPolicyAttachments(&policy)
	if err != nil {
		return nil, err
	}
	policy.AttachmentCount = attachments

	store.CachedPolicies[policy.PolicyID] = &CachedPolicy{
		Policy:   &policy,
		LoadedAt: time.Now(),
	}

	return &policy, nil
}

func (store *Store) GetPolicyByID(policyID string) (*Policy, error) {
	if cachedPolicy, ok := store.CachedPolicies[policyID]; ok {
		if cachedPolicy.NeverExpire || time.Since(cachedPolicy.LoadedAt) < store.TTL {
			return cachedPolicy.Policy, nil
		}

		slog.Debug("invalidating cached policy", "policyID", policyID)
		delete(store.CachedPolicies, policyID)
	}

	var policy Policy

	query := `SELECT * FROM policies WHERE policy_id = $1`
	if err := store.readDB.Get(&policy, query, policyID); err != nil {
		return nil, err
	}

	attachments, err := store.countPolicyAttachments(&policy)
	if err != nil {
		return nil, err
	}
	policy.AttachmentCount = attachments

	store.CachedPolicies[policyID] = &CachedPolicy{
		Policy:   &policy,
		LoadedAt: time.Now(),
	}

	return &policy, nil
}

func (store *Store) getPoliciesByEntityID(arnType ArnType, entityID string) ([]*Policy, error) {
	var policies []*Policy

	var tableName string
	var idFieldName string

	switch arnType {
	case ArnUser:
		tableName = "user_policies"
		idFieldName = "user_id"
	case ArnGroup:
		tableName = "group_policies"
		idFieldName = "group_id"
	default:
		return nil, errors.New("unsupported arn type")
	}

	query_tmpl := `
	SELECT * FROM policies WHERE policy_id = (
		SELECT policy_id FROM %s WHERE %s = $1
	)
	`
	query := fmt.Sprintf(query_tmpl, tableName, idFieldName)

	if err := store.readDB.Select(&policies, query, entityID); err != nil {
		slog.Error("get policies by entity id", "err", err)
		return nil, err
	}

	return policies, nil
}

func (store *Store) countPolicyAttachments(policy *Policy) (int, error) {
	var attachments int
	query := `
	SELECT count(*)
	FROM (
		SELECT 1 FROM user_policies WHERE policy_id = $1
		UNION ALL
		SELECT 1 FROM group_policies WHERE policy_id = $1
	)
	`

	if err := store.readDB.Get(&attachments, query, policy.PolicyID); err != nil {
		return -1, err
	}

	return attachments, nil
}

func GetPolicyHandler(w http.ResponseWriter, r *http.Request) {
	policyArn := r.URL.Query().Get("PolicyArn")
	if policyArn == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("PolicyArn"))
		return
	}

	policy, err := store.GetPolicyByArn(policyArn)
	if err != nil {
		errs.Handle(w, errs.IAMNoSuchEntity(string(errs.ErrEntityTypePolicy), policyArn))
		return
	}

	response := &GetPolicyResponse{
		GetPolicyResult: &GetPolicyResult{
			Policy: policy.Result(),
		},
		ResponseMetadata: &ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}
