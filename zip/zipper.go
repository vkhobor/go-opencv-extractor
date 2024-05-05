package zip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

func Zip(src string, target io.Writer, skipFolders []string) error {
	archive := zip.NewWriter(target)
	defer archive.Close()

	filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = filepath.Join(filepath.Base(src), path[len(src):])

		if info.IsDir() {
			skip := lo.SomeBy(skipFolders, func(s string) bool { return strings.Contains(info.Name(), s) })
			if skip {
				return filepath.SkipDir
			}
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	return nil
}
