package security

import (
	"os"
	"path/filepath"
	"strings"
)

func IsSubdir(baseDir string, path string) bool {
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return false
	}

	absBaseDir = filepath.Clean(absBaseDir) + string(os.PathSeparator)

	absPath, err := filepath.Abs(path)
	if err != nil || !strings.HasPrefix(absPath, absBaseDir) {
		return false
	}

	return true
}
