package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/vkhobor/go-opencv/db_sql"
)

func HandleCreateJob(queries *db_sql.Queries) http.HandlerFunc {
	type jobRequest struct {
		SearchQuery string `json:"search_query"`
		Limit       int64  `json:"limit"`
	}

	type jobResponse struct {
		JobID string `json:"job_id"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			articleRequest := jobRequest{}
			err := render.Decode(r, &articleRequest)
			if err != nil {
				render.Status(r, http.StatusBadRequest)
				render.PlainText(w, r, "Error decoding request body")
				return
			}

			res, err := queries.CreateJob(r.Context(), articleRequest.SearchQuery)
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, "Error creating job")
				return
			}

			render.JSON(w, r, jobResponse{JobID: res})
		},
	)
}
