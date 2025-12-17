package iam

import (
	"database/sql/driver"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"time"
)

const (
	defaultPolicyPath    = "/"
	defaultPolicyVersion = "v1"

	latestPolicyDocumentVersion = "2012-10-17"

	adminPolicy = "admin-policy"
)

type CachedPolicy struct {
	Policy      *Policy
	LoadedAt    time.Time
	NeverExpire bool
}

type Policy struct {
	PolicyID string `db:"policy_id" xml:"PolicyId"`
	Name     string `db:"name" xml:"PolicyName"`
	Arn      string `db:"arn"`

	AttachmentCount int

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Document *PolicyDocument `db:"document" xml:"-"`
}

type PolicyDocument struct {
	Version   string
	Statement []Statement
}

type Statement struct {
	Effect   Effect
	Action   Actions
	Resource Resources
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

func (policy *Policy) Result() *PolicyResult {
	policyPath := "/"

	return &PolicyResult{
		PolicyName:       policy.Name,
		DefaultVersionId: defaultPolicyVersion,
		PolicyId:         policy.PolicyID,
		Path:             policyPath,
		Arn:              policy.Arn,
		AttachmentCount:  policy.AttachmentCount,
		CreateDate:       policy.CreatedAt,
		UpdateDate:       policy.UpdatedAt,
	}
}

func (pd *PolicyDocument) Value() (driver.Value, error) {
	return json.Marshal(pd)
}

func (pd *PolicyDocument) Scan(src any) error {
	if src == nil {
		*pd = PolicyDocument{}
		return nil
	}

	switch s := src.(type) {
	case []byte:
		return json.Unmarshal(s, pd)
	case string:
		return json.Unmarshal([]byte(s), pd)
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}
}
