package handlers

import (
	"strings"

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
	paths, err := h.Client.ListBuckets()
	if err != nil {
		tui.PrintError(err)
	}

	tui.PrintFileList(paths)
}

func (h *S3Handler) ListObjects(path string) {
	parts := strings.Split(path, "/")
	bucket := parts[0]
	key := strings.Join(parts[1:], "/")

	paths, err := h.Client.ListObjects(bucket, key)
	if err != nil {
		tui.PrintError(err)
	}

	tui.PrintFileList(paths)
}
