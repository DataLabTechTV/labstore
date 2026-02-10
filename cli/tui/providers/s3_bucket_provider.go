package providers

import (
	"context"
	"time"

	"github.com/IllumiKnowLabs/labstore/cli/credentials"
	"github.com/IllumiKnowLabs/labstore/cli/errs"
	"github.com/IllumiKnowLabs/labstore/client/s3"
	"github.com/IllumiKnowLabs/labstore/server/config"
)

type S3BucketProvider struct {
	Active  bool
	Host    string
	Port    uint16
	Profile string
}

func (p *S3BucketProvider) Select(args ...string) error {
	if len(args) < 1 {
		return &errs.ErrInsufficientArguments{}
	}

	profile := args[0]
	p.Active = profile != ""
	p.Host = config.App.Server.S3.Address.Host
	p.Port = config.App.Server.S3.Address.Port
	p.Profile = profile
	return nil
}

func (p *S3BucketProvider) Deselect() error {
	p.Active = false
	p.Host = ""
	p.Port = 0
	p.Profile = ""
	return nil
}

func (p *S3BucketProvider) Selected() string {
	return p.Profile
}

func (p *S3BucketProvider) LastSelected() (string, bool) {
	// No state, because there is no navigation.
	return "", false
}

func (p *S3BucketProvider) Children() ([]Entry, error) {
	if !p.Active {
		return []Entry{}, nil
	}

	profile, err := credentials.LoadProfile(p.Profile)
	if err != nil {
		return nil, err
	}

	client := s3.NewS3Client(
		context.Background(),
		p.Host,
		p.Port,
		profile.AccessKey,
		profile.SecretKey,
		false,
		config.App.Server.S3.IO.BufferSize,
	)

	buckets, err := client.ListBuckets()
	if err != nil {
		return nil, err
	}

	var entries []Entry
	for _, bucket := range buckets {
		entry := Entry{
			Name:    bucket.Name,
			ModTime: time.Time(bucket.CreationDate),
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (p *S3BucketProvider) Stat(path string) (Entry, error) {
	return Entry{}, nil
}

func (p *S3BucketProvider) Delete(path string) error {
	return nil
}
