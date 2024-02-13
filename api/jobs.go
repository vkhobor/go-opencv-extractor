package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/vkhobor/go-opencv/db_sql"
	"github.com/vkhobor/go-opencv/jobs"
)

func HandleCreateJob(queries *db_sql.Queries, jobCreator *jobs.JobCreator) http.HandlerFunc {
	type jobRequest struct {
		SearchQuery string `json:"search_query"`
		Limit       int64  `json:"limit"`
	}

	type jobResponse struct {
		Id string `json:"id"`
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

			res, err := queries.CreateJob(r.Context(), db_sql.CreateJobParams{
				SearchQuery: sql.NullString{
					String: articleRequest.SearchQuery,
					Valid:  true,
				},
				Limit: sql.NullInt64{
					Int64: articleRequest.Limit,
					Valid: true,
				},
				ID: uuid.New().String(),
			})
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}
			fmt.Printf("Starting job with id %v\n", res.ID)
			go jobCreator.RunScrapeJob(articleRequest.SearchQuery, int(articleRequest.Limit), res.ID)

			render.JSON(w, r, jobResponse{Id: res.ID})
		},
	)
}

func HandleListJobs(queries *db_sql.Queries) http.HandlerFunc {
	type jobProgress struct {
		Total     int64 `json:"total"`
		Completed int64 `json:"completed"`
	}

	type jobResponse struct {
		ID          string      `json:"id"`
		SearchQuery string      `json:"search_query"`
		Limit       int64       `json:"limit"`
		Progesss    jobProgress `json:"progress"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			res, err := queries.ListJobsWithProgress(r.Context())
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}

			type job struct {
				ID            string
				NumOfImported int64
				All           int64
				SearchQuery   string
				Limit         int64
			}
			jobs := make(map[string]*job)
			for _, j := range res {
				hm := jobs[j.ID]
				if hm == nil {
					jk := &job{
						ID:            j.ID,
						NumOfImported: 0,
						SearchQuery:   j.SearchQuery.String,
						Limit:         j.Limit.Int64,
					}
					jobs[j.ID] = jk
					hm = jk
				}
				hm.All++
				if j.Status.String == "imported" {
					hm.NumOfImported++
				}
			}

			jobsResponse := []jobResponse{}
			for _, j := range jobs {
				jobsResponse = append(jobsResponse, jobResponse{
					SearchQuery: j.SearchQuery,
					ID:          j.ID,
					Limit:       j.Limit,
					Progesss: jobProgress{
						Total:     j.All,
						Completed: j.NumOfImported,
					},
				})
			}

			render.JSON(w, r, jobsResponse)
		},
	)
}
