package providers

import (
	"context"

	"github.com/IllumiKnowLabs/labstore/cli/credentials"
)

type ProfilesProvider struct {
	ActiveProfile *string
}

var initialized bool = false

func (p *ProfilesProvider) Enter(ctx context.Context, path string) error {
	p.ActiveProfile = &path
	return nil
}

func (p *ProfilesProvider) State() (string, bool) {
	return "", false
}

func (*ProfilesProvider) List(ctx context.Context) ([]Entry, error) {
	if !initialized {
		credentials.Init()
	}

	cred, err := credentials.LoadCredentials()
	if err != nil {
		return nil, err
	}

	var entries []Entry
	for profile := range cred.Profiles {
		entries = append(entries, Entry{Name: profile})
	}
	return entries, nil
}

func (*ProfilesProvider) Stat(ctx context.Context, path string) (Entry, error) {
	return Entry{}, nil
}

func (*ProfilesProvider) Delete(ctx context.Context, path string) error {
	return nil
}
