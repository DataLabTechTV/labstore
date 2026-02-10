package providers

type S3FSProvider struct {
	Bucket *string
	Key    *string
}

func (p *S3FSProvider) Select(path string) error {
	return nil
}

func (p *S3FSProvider) Deselect() error {
	return nil
}

func (p *S3FSProvider) Selected() string {
	return ""
}

func (p *S3FSProvider) LastSelected() (string, bool) {
	return "", false
}

func (p *S3FSProvider) Children() ([]Entry, error) {
	return nil, nil
}

func (p *S3FSProvider) Stat(path string) (Entry, error) {
	return Entry{}, nil
}

func (p *S3FSProvider) Delete(path string) error {
	return nil
}
