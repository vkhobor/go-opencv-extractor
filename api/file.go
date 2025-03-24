package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/db"
	"github.com/vkhobor/go-opencv/zip"
)

func ExportWorkspace(config config.DirectoryConfig) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Set the Content-Type header to application/zip
			w.Header().Set("Content-Type", "application/zip")

			// Set the Content-Disposition header so the browser knows it's an attachment
			w.Header().Set("Content-Disposition", "attachment; filename=images.zip")
			zip.ZipFromPath(config.GetImagesDir(), w, []string{"videos", "references"})
		},
	)
}

func HandleFileServeById(queries *db.Queries) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fileIdParam := chi.URLParam(r, "id")
			slog.Debug("Serving file by id", "fileIdParam", fileIdParam)
			if fileIdParam == "" {
				render.Status(r, http.StatusBadRequest)
				render.PlainText(w, r, "No file id provided")
				return
			}
			res, err := queries.GetBlob(r.Context(), fileIdParam)
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}

			http.ServeFile(w, r, res)
		},
	)
}
