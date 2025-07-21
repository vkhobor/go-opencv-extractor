package api

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/vkhobor/go-opencv/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/vkhobor/go-opencv/internal/config"
	"github.com/vkhobor/go-opencv/internal/zip"
)

func ExportWorkspace(config config.DirectoryConfig) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/zip")

			w.Header().Set("Content-Disposition", "attachment; filename=images.zip")
			zip.ZipFromPath(config.GetImagesDir(), w, []string{"videos", "references"})
		},
	)
}

func HandleFileServeById(sqlDB *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fileIdParam := chi.URLParam(r, "id")
			slog.Debug("Serving file by id", "fileIdParam", fileIdParam)
			if fileIdParam == "" {
				render.Status(r, http.StatusBadRequest)
				render.PlainText(w, r, "No file id provided")
				return
			}
			queries := db.New(sqlDB)
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
