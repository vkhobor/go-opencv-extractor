package path

import (
	"os"
	"path/filepath"
	"strings"
)

func MustEnsurePath(path string, isDir bool) {
	err := EnsurePath(path, isDir)
	if err != nil {
		panic(err)
	}
}

func EnsurePath(path string, isDir bool) error {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		path = filepath.Join(home, path[2:])
	}

	if !isDir {
		path = filepath.Dir(path)
	}

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	return nil
}
