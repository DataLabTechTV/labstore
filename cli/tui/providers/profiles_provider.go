package providers

import (
	"github.com/IllumiKnowLabs/labstore/cli/credentials"
)

type ProfilesProvider struct {
	ActiveProfile *string
}

var initialized bool = false

func (p *ProfilesProvider) Select(profile string) error {
	p.ActiveProfile = &profile
	return nil
}

func (p *ProfilesProvider) Deselect() error {
	p.ActiveProfile = nil
	return nil
}

func (p *ProfilesProvider) Selected() string {
	return *p.ActiveProfile
}

func (p *ProfilesProvider) LastSelected() (string, bool) {
	return "", false
}

func (*ProfilesProvider) Children() ([]Entry, error) {
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

func (*ProfilesProvider) Stat(path string) (Entry, error) {
	return Entry{}, nil
}

func (*ProfilesProvider) Delete(path string) error {
	return nil
}
