package iam

import "time"

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
