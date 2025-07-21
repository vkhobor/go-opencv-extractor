package path

import (
	"os"
	"path/filepath"
)

func MustEnsurePath(path string, isDir bool) {
	err := EnsurePath(path, isDir)
	if err != nil {
		panic(err)
	}
}

func EnsurePath(path string, isDir bool) error {
	if !isDir {
		path = filepath.Dir(path)
	}

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	return nil
}
