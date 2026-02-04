package providers

import "context"

type S3BucketProvider struct {
	Bucket *string
}

func (p *S3BucketProvider) Enter(ctx context.Context, path string) error {
	return nil
}

func (p *S3BucketProvider) State() (string, bool) {
	return "", false
}

func (p *S3BucketProvider) List(ctx context.Context) ([]Entry, error) {
	return nil, nil
}

func (p *S3BucketProvider) Stat(ctx context.Context, path string) (Entry, error) {
	return Entry{}, nil
}

func (p *S3BucketProvider) Delete(ctx context.Context, path string) error {
	return nil
}
