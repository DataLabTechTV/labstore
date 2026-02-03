package providers

import (
	"context"
	"time"
)

type Entry struct {
	Name    string
	Path    string
	IsDir   bool
	Size    int64
	ModTime time.Time
}

type Provider interface {
	SetPath(ctx context.Context, path string) error
	List(ctx context.Context) ([]Entry, error)
	Stat(ctx context.Context, path string) (Entry, error)
	Delete(ctx context.Context, path string) error
}

func NewFSProvider() FSProvider {
	return FSProvider{Path: "."}
}

func NewS3Provider() S3Provider {
	// TODO: plugin bucket and key
	return S3Provider{Bucket: "", Key: ""}
}
