package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/vkhobor/go-opencv/db_sql"
)

func HandleAppState(queries *db_sql.Queries) http.HandlerFunc {

	type appStateResponse struct {
		State string `json:"state"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			res, err := queries.GetReferences(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if len(res) == 0 {
				render.JSON(w, r, appStateResponse{State: "no_references"})
				return
			}

			render.JSON(w, r, appStateResponse{State: "ok"})
		},
	)
}
