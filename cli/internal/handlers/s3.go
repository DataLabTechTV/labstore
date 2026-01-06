package handlers

import (
	"github.com/IllumiKnowLabs/labstore/backend/pkg/helper"
	"github.com/IllumiKnowLabs/labstore/cli/internal/tui"
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

func (h *S3Handler) ListBuckets() {
	buckets, err := h.Client.ListBuckets()
	if err != nil {
		tui.PrintError(err)
		return
	}

	tui.PrintBuckets(buckets)
}

func (h *S3Handler) ListObjects(bucket string, key *string) {
	if key == nil {
		key = helper.StringPtr("/")
	}

	res := h.Client.ListObjects(bucket, *key, true)
	for r := range res {
		if r.Err != nil {
			tui.PrintError(r.Err)
			return
		}

		tui.PrintObject(bucket, r.Value)
	}
}
