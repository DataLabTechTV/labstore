package helper

import (
	"os"
	"path/filepath"
	"strings"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsExist(err)
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func TildePath(path string) string {
	home := os.Getenv("HOME")

	if strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}

	return path
}

func MustResolveToRelativePath(path string) string {
	cwd := Must(os.Getwd())
	absPath := Must(filepath.Abs(path))
	relPath := Must(filepath.Rel(cwd, absPath))
	return relPath
}
