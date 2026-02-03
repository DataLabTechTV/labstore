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
	List(ctx context.Context, path string) ([]Entry, error)
	Stat(ctx context.Context, path string) (Entry, error)
	Delete(ctx context.Context, path string) error
}

func NewFSProvider() FSProvider {
	return FSProvider{}
}

func NewS3Provider() S3Provider {
	return S3Provider{}
}
