package providers

import (
	"os"
	"path/filepath"
	"time"

	"github.com/IllumiKnowLabs/labstore/server/helper"
)

type FSProvider struct {
	Path  string
	state map[string]string
}

func (p *FSProvider) Select(path string) error {
	p.state[p.Path] = path
	p.Path = filepath.Join(p.Path, path)
	return nil
}

func (p *FSProvider) Deselect() error {
	return nil
}

func (p *FSProvider) LastSelected() (string, bool) {
	state, ok := p.state[p.Path]
	return state, ok
}

func (p *FSProvider) Selected() string {
	homeRelPath := helper.TildePath(p.Path)
	return homeRelPath
}

func (p *FSProvider) Children() ([]Entry, error) {
	entries := []Entry{}

	if parent, ok := parentDir(p.Path); ok {
		var size int64
		var modTime time.Time

		info, err := os.Stat(parent)
		if err != nil {
			size = 0
			modTime = time.Unix(0, 0)
		} else {
			size = info.Size()
			modTime = info.ModTime()
		}

		entry := Entry{
			Name:    "..",
			Path:    parent,
			IsDir:   true,
			Size:    size,
			ModTime: modTime,
		}

		entries = append(entries, entry)
	}

	dirEntries, err := os.ReadDir(p.Path)
	if err != nil {
		return nil, err
	}

	for _, dirEntry := range dirEntries {
		info, err := dirEntry.Info()
		if err != nil {
			return nil, err
		}

		entry := Entry{
			Name:    dirEntry.Name(),
			Path:    filepath.Join(p.Path, dirEntry.Name()),
			IsDir:   dirEntry.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (p *FSProvider) Stat(path string) (Entry, error) {
	return Entry{}, nil
}

func (p *FSProvider) Delete(path string) error {
	return nil
}

func parentDir(path string) (string, bool) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", false
	}

	parent := filepath.Dir(abs)

	if parent != abs {
		return parent, true
	}

	return "", false
}
