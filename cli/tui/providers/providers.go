package providers

import (
	"context"
	"time"

	"github.com/IllumiKnowLabs/labstore/server/helper"
)

type Entry struct {
	Name    string
	Path    string
	IsDir   bool
	Size    int64
	ModTime time.Time
}

type Provider interface {
	Enter(ctx context.Context, path string) error
	List(ctx context.Context) ([]Entry, error)
	Stat(ctx context.Context, path string) (Entry, error)
	Delete(ctx context.Context, path string) error
}

func NewFSProvider() FSProvider {
	return FSProvider{Path: "."}
}

func NewS3FSProvider() S3FSProvider {
	// TODO: plugin bucket and key
	return S3FSProvider{Bucket: NewS3BucketProvider().Bucket, Key: helper.Ptr("")}
}

func NewS3BucketProvider() S3BucketProvider {
	// TODO: plugin bucket
	return S3BucketProvider{Bucket: helper.Ptr("")}
}
