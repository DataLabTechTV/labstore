package providers

import (
	"github.com/IllumiKnowLabs/labstore/cli/credentials"
	"github.com/IllumiKnowLabs/labstore/cli/errs"
)

type ProfilesProvider struct {
	ActiveProfile *string
}

var initialized bool = false

func NewProfilesProvider() *ProfilesProvider {
	// TODO: keep state of last selected profile (fallback to default if it does not exist anymore)
	return &ProfilesProvider{ActiveProfile: nil}
}

func (p *ProfilesProvider) Select(args ...string) error {
	if len(args) < 1 {
		return &errs.ErrInsufficientArguments{}
	}

	p.ActiveProfile = &args[0]
	return nil
}

func (p *ProfilesProvider) Deselect() error {
	p.ActiveProfile = nil
	return nil
}

func (p *ProfilesProvider) Selected() string {
	return *p.ActiveProfile
}

func (*ProfilesProvider) Children() ([]Entry, error) {
	if !initialized {
		credentials.Init()
		initialized = true
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

func (*ProfilesProvider) Delete(path string) error {
	return nil
}
