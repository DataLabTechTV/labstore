package providers

import (
	"path/filepath"
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
	Enter(string) error
	State() (string, bool)
	CWD() string
	List() ([]Entry, error)
	Stat(string) (Entry, error)
	Delete(string) error
}

func NewS3BucketProvider() *S3BucketProvider {
	// TODO: keep state of last selected bucket (might not exist anymore)
	return &S3BucketProvider{}
}

func NewProfilesProvider() *ProfilesProvider {
	return &ProfilesProvider{ActiveProfile: nil}
}

func NewS3FSProvider() *S3FSProvider {
	// TODO: keep state of last selected bucket and key (both might not exist anymore)
	return &S3FSProvider{Bucket: nil, Key: nil}
}

func NewFSProvider(rootPath string) *FSProvider {
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		absPath = rootPath
	}

	return &FSProvider{
		Path:  absPath,
		state: make(map[string]string),
	}
}

func EntryNames(entries []Entry) []string {
	names := make([]string, len(entries))
	for i, entry := range entries {
		names[i] = entry.Name
	}
	return names
}
