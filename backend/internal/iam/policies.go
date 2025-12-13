package iam

import (
	"encoding/xml"
	"log/slog"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
)

const adminPolicy = "admin-policy"
const latestPolicyDocumentVersion = "2012-10-17"

type Policy struct {
	PolicyID  string    `db:"policy_id"`
	Name      string    `db:"name"`
	CreateAt  time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Document  *PolicyDocument
}

type PolicyDocument struct {
	Version   string
	Statement []Statement
}

type Statement struct {
	Effect    Effect
	Actions   []Action
	Resources []string
}

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
	AttachmentCount  uint
	CreateDate       time.Time
	UpdateDate       time.Time
}

func GetPolicyByID(policyID string) (*Policy, error) {
	// TODO: implement
	return nil, nil
}

func CreatePolicy(name string, doc *PolicyDocument) (*Policy, error) {
	policyID := GenerateUniqueID(ManagedPolicyUniqueID)

	policy := &Policy{
		PolicyID: policyID,
		Name:     name,
		Document: doc,
	}

	query := `
	INSERT INTO policies (policy_id, name, document)
	VALUES (:policy_id, :name, :document)
	`

	_, err := store.writeDB.NamedExec(query, &policy)
	if err != nil {
		slog.Error("create policy insert", "err", err)
		return nil, err
	}

	policy, err = GetPolicyByID(policyID)
	if err != nil {
		slog.Error("get policy by id", "err", err)
		return nil, &errs.ErrNotFound{Type: errs.ErrEntityTypePolicy, Resource: policyID}
	}

	return policy, nil

}
