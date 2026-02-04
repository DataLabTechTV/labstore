package providers

import "context"

type S3BucketProvider struct {
	Bucket *string
}

func (S3BucketProvider) Enter(ctx context.Context, path string) error {
	return nil
}

func (S3BucketProvider) List(ctx context.Context) ([]Entry, error) {
	return nil, nil
}

func (S3BucketProvider) Stat(ctx context.Context, path string) (Entry, error) {
	return Entry{}, nil
}

func (S3BucketProvider) Delete(ctx context.Context, path string) error {
	return nil
}
