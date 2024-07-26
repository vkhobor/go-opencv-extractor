package filters

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/db"
)

type filter struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func HandleGetFilters(queries *db.Queries) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			res, err := queries.GetFilters(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if len(res) == 0 {
				render.JSON(w, r, []filter{})
				return
			}

			filters := lo.Map(res, func(item db.GetFiltersRow, i int) filter {
				return filter{
					Name: item.Name.String,
					ID:   item.ID,
				}
			})

			render.JSON(w, r, filters)
		},
	)
}
