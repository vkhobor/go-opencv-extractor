package path

import (
	"os"
	"path/filepath"
	"strings"
)

func EnsurePath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		path = filepath.Join(home, path[2:])
	}

	dirPath := filepath.Dir(path)
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return "", err
	}
	return path, nil
}

func IsDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}
