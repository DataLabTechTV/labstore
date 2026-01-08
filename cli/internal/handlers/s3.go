package handlers

import (
	"fmt"
	"os"

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
	tui.PrintTitle("ListBuckets")

	buckets, err := h.Client.ListBuckets()
	if err != nil {
		tui.PrintError(err)
		return
	}

	for _, bucket := range buckets {
		tui.PrintBucket(bucket)
	}
}

func (h *S3Handler) ListObjects(bucket string, key *string) {
	if key == nil {
		key = helper.StringPtr("/")
	}

	tui.PrintTitle(fmt.Sprintf("ListObjectsV2: %s", bucket))

	res := h.Client.ListObjects(bucket, *key, true)
	for r := range res {
		if r.Err != nil {
			tui.PrintError(r.Err)
			return
		}

		switch {
		case r.IsObject():
			tui.PrintObject(*r.Object)
		case r.IsCommonPrefix():
			tui.PrintCommonPrefix(*r.CommonPrefix)
		}
	}
}

func (h *S3Handler) PutObject(bucket, key, localPath string) {
	if !helper.FileExists(localPath) {
		tui.PrintError(fmt.Errorf("file not found: %s", localPath))
		return
	}

	tui.PrintTitle(fmt.Sprintf("PutObject: %s/%s", bucket, key))

	file, err := os.Open(localPath)
	if err != nil {
		tui.PrintError(err)
		return
	}

	err = h.Client.PutObject(bucket, key, file, tui.NewProgressBar())
	if err != nil {
		tui.PrintError(err)
		return
	}

	tui.PrintSuccess(fmt.Sprintf("%s uploaded", localPath))
}
