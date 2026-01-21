package iam

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/IllumiKnowLabs/labstore/server/internal/core"
	"github.com/IllumiKnowLabs/labstore/server/pkg/errs"
	"github.com/IllumiKnowLabs/labstore/server/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/server/pkg/types"
)

func (store *Store) CreatePolicy(ctx context.Context, name string, doc *PolicyDocument) (*Policy, error) {
	policyID := GenerateUniqueID(ManagedPolicyUniqueID)

	policy := &Policy{
		PolicyID: policyID,
		Name:     name,
		Arn:      toArn(ArnPolicy, defaultPolicyPath+name),
		Document: doc,
	}

	query := `
	INSERT INTO policies (policy_id, name, arn, document)
	VALUES (:policy_id, :name, :arn, :document)
	`

	_, err := store.sqlNamedExecContext(ctx, query, &policy)
	if err != nil {
		slog.Warn("create policy", "err", err)
		return nil, &errs.ErrExists{Type: errs.ErrEntityTypePolicy, Resource: policyID}
	}

	policy, err = store.GetPolicyByID(ctx, policyID)
	if err != nil {
		slog.Error("get policy by id", "err", err)
		return nil, &errs.ErrNotFound{Type: errs.ErrEntityTypePolicy, Resource: policyID}
	}

	store.CachedPolicies[policyID] = &CachedPolicy{
		Policy:   policy,
		LoadedAt: time.Now(),
	}

	return policy, nil
}

func CreatePolicyHandler(w http.ResponseWriter, r *http.Request) {
	name := r.Form.Get("PolicyName")
	if name == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("PolicyName"))
		return
	}

	document := r.Form.Get("PolicyDocument")
	if document == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("PolicyDocument"))
		return
	}

	var doc PolicyDocument
	err := json.Unmarshal([]byte(document), &doc)
	if err != nil {
		errs.Handle(w, err)
		return
	}

	ctx := r.Context()

	policy, err := store.CreatePolicy(ctx, name, &doc)
	if err != nil {
		var errNotFound *errs.ErrNotFound
		if errors.As(err, &errNotFound) {
			errs.Handle(w, errs.IAMServiceFailure())
			return
		}

		var errExists *errs.ErrExists
		if errors.As(err, &errExists) {
			errs.Handle(w, errs.IAMEntityAlreadyExists(errExists.Resource))
			return
		}

		errs.Handle(w, errs.IAMServiceFailure())
		return
	}

	response := &types.CreatePolicyResponse{
		CreatePolicyResult: &types.CreatePolicyResult{
			Policy: policy.Result(),
		},
		ResponseMetadata: &types.ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	helper.WriteXMLResponse(w, http.StatusOK, response)
}
