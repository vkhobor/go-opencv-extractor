package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/db"
)

var dummyFilterId = "1fed33d4-0ea3-4b84-909c-261e4b2a3d43"

type filterImage struct {
	BlobId string `json:"blob_id"`
}

type filter struct {
	Name         string        `json:"name"`
	ID           string        `json:"id"`
	FilterImages []filterImage `json:"filter_images"`
}

func HandleGetFilters(queries *db.Queries) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			res, err := queries.GetReferences(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if len(res) == 0 {
				render.JSON(w, r, []filter{})
				return
			}

			singleFilter := filter{
				ID:   dummyFilterId,
				Name: "Default Filter",
				FilterImages: lo.Map(res, func(item db.GetReferencesRow, index int) filterImage {
					return filterImage{
						BlobId: item.BlobStorageID,
					}
				}),
			}

			render.JSON(w, r, []filter{singleFilter})
		},
	)
}

func HandleGetFilter(queries *db.Queries) http.HandlerFunc {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			if id == "" {
				render.Status(r, http.StatusBadRequest)
				render.PlainText(w, r, "No id provided")
				return
			}

			res, err := queries.GetReferences(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			response := filter{
				ID:   dummyFilterId,
				Name: "Default Filter",
				FilterImages: lo.Map(res, func(item db.GetReferencesRow, index int) filterImage {
					return filterImage{
						BlobId: item.BlobStorageID,
					}
				}),
			}

			render.JSON(w, r, response)
		},
	)
}
