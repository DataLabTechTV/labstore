package iam

import (
	"database/sql/driver"
	"encoding/json"
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
	policy      *Policy
	loadedAt    time.Time
	neverExpire bool
}

type Policy struct {
	PolicyID string `db:"policy_id"`
	Name     string `db:"name"`
	Arn      string `db:"arn"`

	Document *PolicyDocument `db:"document"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	AttachmentCount int
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
