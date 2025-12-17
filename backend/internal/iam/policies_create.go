package iam

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/internal/core"
	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
)

type CreatePolicyResponse struct {
	XMLName            xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ CreatePolicyResponse"`
	CreatePolicyResult *CreatePolicyResult
	ResponseMetadata   *ResponseMetadata
}

type CreatePolicyResult struct {
	Policy *PolicyResult
}

type PolicyResult struct {
	XMLName          xml.Name `xml:"Policy"`
	PolicyName       string
	DefaultVersionId string
	PolicyId         string
	Path             string
	Arn              string
	AttachmentCount  int
	CreateDate       time.Time
	UpdateDate       time.Time
}

func (store *Store) CreatePolicy(name string, doc *PolicyDocument) (*Policy, error) {
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

	_, err := store.writeDB.NamedExec(query, &policy)
	if err != nil {
		slog.Warn("create policy", "err", err)
		return nil, &errs.ErrExists{Type: errs.ErrEntityTypePolicy, Resource: policyID}
	}

	policy, err = store.GetPolicyByID(policyID)
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
	name := r.URL.Query().Get("PolicyName")
	if name == "" {
		errs.Handle(w, errs.HTTPMissingQueryParam("PolicyName"))
		return
	}

	document := r.URL.Query().Get("PolicyDocument")
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

	policy, err := store.CreatePolicy(name, &doc)
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

	policyPath := "/"

	response := &CreatePolicyResponse{
		CreatePolicyResult: &CreatePolicyResult{
			Policy: &PolicyResult{
				PolicyName:       policy.Name,
				DefaultVersionId: defaultPolicyVersion,
				PolicyId:         policy.PolicyID,
				Path:             policyPath,
				Arn:              policy.Arn,
				AttachmentCount:  policy.AttachmentCount,
				CreateDate:       policy.CreatedAt,
				UpdateDate:       policy.UpdatedAt,
			},
		},
		ResponseMetadata: &ResponseMetadata{
			RequestId: core.NewRequestID(),
		},
	}

	core.WriteXML(w, http.StatusOK, response)
}
