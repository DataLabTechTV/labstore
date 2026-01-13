package handlers

import (
	"github.com/IllumiKnowLabs/labstore/client/s3"
)

type S3Handler struct {
	Client *s3.S3Client
}

func NewS3Handler(client *s3.S3Client) *S3Handler {
	return &S3Handler{
		Client: client,
	}
}
