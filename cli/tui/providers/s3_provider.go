package providers

import "context"

type S3Provider struct{}

func (S3Provider) List(ctx context.Context, path string) ([]Entry, error) {
	return nil, nil
}

func (S3Provider) Stat(ctx context.Context, path string) (Entry, error) {
	return Entry{}, nil
}

func (S3Provider) Delete(ctx context.Context, path string) error {
	return nil
}
