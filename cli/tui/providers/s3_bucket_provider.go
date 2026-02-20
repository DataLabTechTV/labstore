package providers

import (
	"context"
	"time"

	"github.com/IllumiKnowLabs/labstore/cli/credentials"
	"github.com/IllumiKnowLabs/labstore/cli/errs"
	"github.com/IllumiKnowLabs/labstore/cli/types"
	"github.com/IllumiKnowLabs/labstore/client/s3"
	"github.com/IllumiKnowLabs/labstore/server/config"
)

type S3BucketProvider struct {
	Active  bool
	Host    string
	Port    uint16
	Profile string
}

func NewS3BucketProvider() *S3BucketProvider {
	// TODO: keep state of last selected bucket (might not exist anymore)
	return &S3BucketProvider{
		Host: config.App.Server.S3.Address.Host,
		Port: config.App.Server.S3.Address.Port,
	}
}

func (p *S3BucketProvider) Select(args ...string) error {
	if len(args) < 1 {
		return &errs.ErrInsufficientArguments{}
	}

	profile := args[0]
	p.Active = profile != ""
	p.Profile = profile
	return nil
}

func (p *S3BucketProvider) Deselect() {
	p.Active = false
	p.Profile = ""
}

func (p *S3BucketProvider) Selected() string {
	return p.Profile
}

func (p *S3BucketProvider) Children() ([]types.Entry, error) {
	if !p.Active {
		return []types.Entry{}, nil
	}

	s3Client, err := p.newClient()
	if err != nil {
		return nil, err
	}

	buckets, err := s3Client.ListBuckets()
	if err != nil {
		return nil, err
	}

	var entries []types.Entry
	for _, bucket := range buckets {
		entry := types.Entry{
			Name:    bucket.Name,
			ModTime: time.Time(bucket.CreationDate),
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (p *S3BucketProvider) Stat(bucket string) (*types.Entry, error) {
	return &types.Entry{}, nil
}

func (p *S3BucketProvider) Delete(bucket string) error {
	if !p.Active {
		return nil
	}

	s3Client, err := p.newClient()
	if err != nil {
		return err
	}

	_, err = s3Client.DeleteBucket(bucket)
	if err != nil {
		return err
	}

	return nil
}

func (p *S3BucketProvider) newClient() (*s3.Client, error) {
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

	return client, nil
}
