package handlers

import (
	"github.com/IllumiKnowLabs/labstore/client/pkg/s3"
)

type S3Handler struct {
	Client *s3.Client
}

func NewS3Handler(client *s3.Client) *S3Handler {
	return &S3Handler{
		Client: client,
	}
}
