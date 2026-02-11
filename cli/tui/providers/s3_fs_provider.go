package providers

import (
	"context"
	"fmt"
	"path"
	"strings"
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

	lastSelected map[string]string
}

func (p *S3FSProvider) Select(args ...string) error {
	if len(args) < 2 {
		return &errs.ErrInsufficientArguments{}
	}

	p.Active = true
	p.Profile = args[0]
	p.Bucket = args[1]

	lastSelKey := p.lastSelKey()

	if len(args) > 2 {
		p.Key = path.Clean(args[2])
		if p.Key == "." {
			p.Key = ""
		} else {
			p.Key += "/"
		}
	}

	p.lastSelected[lastSelKey] = p.Key

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
	return p.Key
}

func (p *S3FSProvider) LastSelected() (string, bool) {
	state, ok := p.lastSelected[p.lastSelKey()]
	return state, ok

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

	if keyParts := strings.Split(p.Key, "/"); len(keyParts) > 1 {
		upLevelEntry := Entry{
			Name:    "..",
			Path:    path.Clean(keyParts[len(keyParts)-2]) + "/",
			IsDir:   true,
			ModTime: time.Now(),
		}
		entries = append(entries, upLevelEntry)
	}

	for res := range resCh {
		if res.Err != nil {
			return nil, res.Err
		}

		var entry *Entry
		if res.CommonPrefix != nil {
			entry = &Entry{
				Name:    path.Base(res.CommonPrefix.Prefix) + "/",
				Path:    res.CommonPrefix.Prefix,
				IsDir:   true,
				ModTime: time.Now(),
			}
		} else {
			entry = &Entry{
				Name:    path.Base(res.Object.Key),
				Path:    res.Object.Key,
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

func (p *S3FSProvider) lastSelKey() string {
	return fmt.Sprintf("%s@%s/%s", p.Profile, p.Bucket, p.Key)
}
