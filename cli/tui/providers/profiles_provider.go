package providers

import (
	"github.com/IllumiKnowLabs/labstore/cli/credentials"
	"github.com/IllumiKnowLabs/labstore/cli/types"
)

type ProfilesProvider struct{}

var initialized bool = false

func NewProfilesProvider() *ProfilesProvider {
	// TODO: keep state of last selected profile (fallback to default if it does not exist anymore)
	return &ProfilesProvider{}
}

func (*ProfilesProvider) Children() ([]types.Entry, error) {
	if !initialized {
		credentials.Init()
		initialized = true
	}

	cred, err := credentials.LoadCredentials()
	if err != nil {
		return nil, err
	}

	var entries []types.Entry
	for profile := range cred.Profiles {
		entries = append(entries, types.Entry{Name: profile})
	}
	return entries, nil
}

func (*ProfilesProvider) Delete(path string) error {
	return nil
}
