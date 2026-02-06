package providers

import (
	"github.com/IllumiKnowLabs/labstore/cli/credentials"
)

type ProfilesProvider struct {
	ActiveProfile *string
}

var initialized bool = false

func (p *ProfilesProvider) Enter(path string) error {
	p.ActiveProfile = &path
	return nil
}

func (p *ProfilesProvider) State() (string, bool) {
	return "", false
}

func (p *ProfilesProvider) CWD() string {
	return *p.ActiveProfile
}

func (*ProfilesProvider) List() ([]Entry, error) {
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
