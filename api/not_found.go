package api

import (
	"errors"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/vkhobor/go-opencv/web"
)

func HandleCatchAll() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			prefixedPath := filepath.Join(web.PrefixForClientFiles, r.URL.Path)
			err := readFileToResponse(web.StaticFiles, prefixedPath, w)
			if err == nil {
				return
			}

			err = readFileToResponse(web.StaticFiles, web.IndexHtml, w)
			if err != nil {
				panic(err)
			}
		})
}

var ErrDir = errors.New("path is dir")

func readFileToResponse(fileSystem fs.FS, requestedPath string, w http.ResponseWriter) error {
	f, err := fileSystem.Open(requestedPath)
	if err != nil {
		return err
	}
	defer f.Close()

	stat, _ := f.Stat()
	if stat.IsDir() {
		return ErrDir
	}

	contentType := mime.TypeByExtension(filepath.Ext(requestedPath))
	w.Header().Set("Content-Type", contentType)
	_, err = io.Copy(w, f)
	return err
}
