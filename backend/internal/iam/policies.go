package iam

const adminPolicy = "admin-policy"
const latestPolicyDocumentVersion = "2012-10-17"

type Policy struct {
	ID       string
	Document *PolicyDocument
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
