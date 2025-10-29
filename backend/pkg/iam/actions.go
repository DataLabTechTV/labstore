package iam

type Action string

const (
	ListAllMyBuckets Action = "s3:ListAllMyBuckets"
	CreateBucket     Action = "s3:CreateBucket"
	DeleteBucket     Action = "s3:DeleteBucket"
	ListBucket       Action = "s3:ListBucket"
	PutObject        Action = "s3:PutObject"
	GetObject        Action = "s3:GetObject"
	DeleteObject     Action = "s3:DeleteObject"
)
