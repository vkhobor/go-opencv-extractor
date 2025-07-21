package filters

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/render"
	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/internal/db"
)

func HandleGetReferences(sqlDB *sql.DB) http.HandlerFunc {
	type reference struct {
		ID string `json:"id"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			queries := db.New(sqlDB)
			res, err := queries.GetFilters(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			response := lo.FilterMap(res, func(item db.GetFiltersRow, index int) (reference, bool) {
				return reference{
					ID: item.BlobID.String,
				}, item.BlobID.Valid
			})

			render.JSON(w, r, response)
		},
	)
}
