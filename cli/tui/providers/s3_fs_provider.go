package providers

import "context"

type S3FSProvider struct {
	Bucket *string
	Key    *string
}

func (S3FSProvider) Enter(ctx context.Context, path string) error {
	return nil
}

func (S3FSProvider) List(ctx context.Context) ([]Entry, error) {
	return nil, nil
}

func (S3FSProvider) Stat(ctx context.Context, path string) (Entry, error) {
	return Entry{}, nil
}

func (S3FSProvider) Delete(ctx context.Context, path string) error {
	return nil
}
