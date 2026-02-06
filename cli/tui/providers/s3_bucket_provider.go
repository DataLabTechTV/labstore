package providers

type S3BucketProvider struct {
	Bucket *string
}

func (p *S3BucketProvider) Enter(path string) error {
	return nil
}

func (p *S3BucketProvider) State() (string, bool) {
	return "", false
}

func (p *S3BucketProvider) CWD() string {
	return ""
}

func (p *S3BucketProvider) List() ([]Entry, error) {
	return nil, nil
}

func (p *S3BucketProvider) Stat(path string) (Entry, error) {
	return Entry{}, nil
}

func (p *S3BucketProvider) Delete(path string) error {
	return nil
}
