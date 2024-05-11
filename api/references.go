package api

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/db"
)

const (
	megabyte = 1 << 20 // 1 megabyte = 2^20 bytes
)

func HandleReferenceUpload(queries *db.Queries, config config.DirectoryConfig) http.HandlerFunc {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseMultipartForm(10 * megabyte)
			files := r.MultipartForm.File
			for _, headers := range files {
				for _, header := range headers {
					file, err := header.Open()
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					defer file.Close()

					err = os.MkdirAll(config.GetReferencesDir(), os.ModePerm)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					dst, err := os.Create(fmt.Sprintf("%s/%s", config.GetReferencesDir(), header.Filename))
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					defer dst.Close()

					if _, err := io.Copy(dst, file); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					path := fmt.Sprintf("%s/%s", config.GetReferencesDir(), header.Filename)
					id := uuid.NewString()
					err = queries.AddReference(r.Context(), id)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					_, err = queries.AddBlob(r.Context(), db.AddBlobParams{
						ID:   id,
						Path: path,
					})
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
			}
			render.Status(r, http.StatusCreated)
		},
	)
}

func HandleGetReferences(queries *db.Queries) http.HandlerFunc {
	type reference struct {
		ID string `json:"id"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			res, err := queries.GetReferences(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			response := lo.Map(res, func(item db.GetReferencesRow, index int) reference {
				return reference{
					ID: item.BlobStorageID,
				}
			})

			render.JSON(w, r, response)
		},
	)
}

func HandleDeleteAllReferences(queries *db.Queries) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := queries.DeleteReferences(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			render.Status(r, http.StatusNoContent)
		},
	)
}
