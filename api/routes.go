package api

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/vkhobor/go-opencv/db_sql"
	"github.com/vkhobor/go-opencv/frontend"
	"github.com/vkhobor/go-opencv/jobs"
)

func NewRouter(
	queries *db_sql.Queries,
	jobCreator *jobs.JobCreator,
) chi.Router {
	router := chi.NewRouter()
	router.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(60 * time.Second)
		w.Write([]byte(fmt.Sprintf("all done.\n")))
	})
	router.Post("/jobs", HandleCreateJob(queries, jobCreator))
	router.Get("/jobs", HandleListJobs(queries))

	router.NotFound(HAndleNNotFound())
	return router
}

func HAndleNNotFound() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := tryRead(frontend.StaticFiles, "dist/extractor/browser", r.URL.Path, w)
			if err == nil {
				return
			}
			err = tryRead(frontend.StaticFiles, "dist/extractor/browser", "index.html", w)
			if err != nil {
				panic(err)
			}
		})
}

var httpFS = http.FileServer(http.FS(frontend.StaticFiles))

var ErrDir = errors.New("path is dir")

func tryRead(fs embed.FS, prefix, requestedPath string, w http.ResponseWriter) error {
	f, err := fs.Open(path.Join(prefix, requestedPath))
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
