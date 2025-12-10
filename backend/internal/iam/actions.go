package iam

import "github.com/gobwas/glob"

type Action string

const (
	S3ListAllMyBuckets Action = "s3:ListAllMyBuckets"
	S3CreateBucket     Action = "s3:CreateBucket"
	S3DeleteBucket     Action = "s3:DeleteBucket"
	S3ListBucket       Action = "s3:ListBucket"
	S3PutObject        Action = "s3:PutObject"
	S3GetObject        Action = "s3:GetObject"
	S3DeleteObject     Action = "s3:DeleteObject"
)

func matchAction(action Action, stmtActions []Action) bool {
	for _, stmtAction := range stmtActions {
		g := glob.MustCompile(string(stmtAction))
		if g.Match(string(action)) {
			return true
		}
	}

	return false

}
