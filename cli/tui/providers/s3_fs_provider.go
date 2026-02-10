package providers

import (
	"context"
	"fmt"
	"time"

	"github.com/IllumiKnowLabs/labstore/cli/credentials"
	"github.com/IllumiKnowLabs/labstore/cli/errs"
	"github.com/IllumiKnowLabs/labstore/client/s3"
	"github.com/IllumiKnowLabs/labstore/server/config"
)

type S3FSProvider struct {
	Active  bool
	Host    string
	Port    uint16
	Profile string
	Bucket  string
	Key     string
}

func (p *S3FSProvider) Select(args ...string) error {
	if len(args) < 2 {
		return &errs.ErrInsufficientArguments{}
	}

	p.Active = true
	p.Host = config.App.Server.S3.Address.Host
	p.Port = config.App.Server.S3.Address.Port
	p.Profile = args[0]
	p.Bucket = args[1]

	if len(args) > 2 {
		p.Key = args[2]
	}

	return nil
}

func (p *S3FSProvider) Deselect() error {
	p.Active = false
	p.Profile = ""
	p.Bucket = ""
	p.Key = ""
	return nil
}

func (p *S3FSProvider) Selected() string {
	return fmt.Sprintf("s3://%s%s", p.Bucket, p.Key)
}

func (p *S3FSProvider) LastSelected() (string, bool) {
	return "", false
}

func (p *S3FSProvider) Children() ([]Entry, error) {
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

	resCh := client.ListObjects(p.Bucket, p.Key, true)
	var entries []Entry

	for res := range resCh {
		if res.Err != nil {
			return nil, res.Err
		}

		var entry *Entry
		if res.CommonPrefix != nil {
			entry = &Entry{
				Name:  res.CommonPrefix.Prefix,
				IsDir: true,
			}
		} else {
			entry = &Entry{
				Name:    res.Object.Key,
				IsDir:   false,
				ModTime: time.Time(res.Object.LastModified),
				Size:    res.Object.Size,
			}
		}
		entries = append(entries, *entry)
	}

	return entries, nil
}

func (p *S3FSProvider) Stat(path string) (Entry, error) {
	return Entry{}, nil
}

func (p *S3FSProvider) Delete(path string) error {
	return nil
}
