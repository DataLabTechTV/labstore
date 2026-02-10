package providers

import (
	"os"
	"time"
)

type Entry struct {
	Name    string
	Path    string
	IsDir   bool
	Size    int64
	ModTime time.Time
}

type Provider interface {
	Select(...string) error
	Deselect() error
	LastSelected() (string, bool)
	Selected() string
	Children() ([]Entry, error)
	Stat(string) (Entry, error)
	Delete(string) error
}

func NewS3BucketProvider() *S3BucketProvider {
	// TODO: keep state of last selected bucket (might not exist anymore)
	return &S3BucketProvider{}
}

func NewProfilesProvider() *ProfilesProvider {
	// TODO: keep state of last selected profile (fallback to default if it does not exist anymore)
	return &ProfilesProvider{ActiveProfile: nil}
}

func NewS3FSProvider() *S3FSProvider {
	// TODO: keep state of last selected bucket and key (both might not exist anymore)
	return &S3FSProvider{}
}

func NewFSProvider() *FSProvider {
	path, err := os.Getwd()
	if err != nil {
		path = "."
	}

	return &FSProvider{
		Path:         path,
		lastSelected: make(map[string]string),
	}
}

func EntryNames(entries []Entry) []string {
	names := make([]string, len(entries))
	for i, entry := range entries {
		names[i] = entry.Name
	}
	return names
}
