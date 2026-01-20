package iam

import (
	"encoding/json"

	"github.com/gobwas/glob"
)

const (
	S3ListAllMyBuckets Action = "s3:ListAllMyBuckets"
	S3CreateBucket     Action = "s3:CreateBucket"
	S3DeleteBucket     Action = "s3:DeleteBucket"
	S3ListBucket       Action = "s3:ListBucket"
	S3PutObject        Action = "s3:PutObject"
	S3GetObject        Action = "s3:GetObject"
	S3DeleteObject     Action = "s3:DeleteObject"
)

type Action string
type Actions []Action

func (a *Actions) UnmarshalJSON(data []byte) error {
	var single Action
	if err := json.Unmarshal(data, &single); err == nil {
		*a = []Action{single}
		return nil
	}

	var multi []Action
	if err := json.Unmarshal(data, &multi); err != nil {
		return err
	}
	*a = multi
	return nil
}

func matchAction(action Action, stmtActions []Action) bool {
	for _, stmtAction := range stmtActions {
		g := glob.MustCompile(string(stmtAction))
		if g.Match(string(action)) {
			return true
		}
	}

	return false

}
