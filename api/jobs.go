package api

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/db_sql"
	"github.com/vkhobor/go-opencv/jobs"
)

func HandleCreateJob(queries *db_sql.Queries, jobCreator *jobs.JobCreator) http.HandlerFunc {
	type jobRequest struct {
		SearchQuery string `json:"search_query"`
		Limit       int    `json:"limit"`
	}

	type jobResponse struct {
		Id string `json:"id"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			articleRequest := jobRequest{}
			err := render.Decode(r, &articleRequest)
			if err != nil || articleRequest.Limit < 1 || articleRequest.SearchQuery == "" {
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
					Int64: int64(articleRequest.Limit),
					Valid: true,
				},
				ID: uuid.New().String(),
			})
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}

			render.JSON(w, r, jobResponse{Id: res.ID})
		},
	)
}

func HandleListJobs(queries *db_sql.Queries) http.HandlerFunc {
	type jobProgress struct {
		Imported   int `json:"imported"`
		Scraped    int `json:"scraped"`
		Downloaded int `json:"downloaded"`
	}

	type jobResponse struct {
		ID          string      `json:"id"`
		SearchQuery string      `json:"search_query"`
		Limit       int         `json:"limit"`
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
				NumOfImported int
				All           int
				SearchQuery   string
				Limit         int
			}

			grouped := lo.GroupBy(res, func(row db_sql.ListJobsWithProgressRow) string {
				return row.ID
			})

			jobsResponse := []jobResponse{}
			for key, value := range grouped {

				if len(value) == 1 && value[0].ID_2.Valid == false {
					jobsResponse = append(jobsResponse, jobResponse{
						SearchQuery: value[0].SearchQuery.String,
						ID:          key,
						Limit:       int(value[0].Limit.Int64),
						Progesss: jobProgress{
							Imported:   0,
							Scraped:    0,
							Downloaded: 0,
						},
					})
					continue
				}

				scraped := len(value)

				downloaded := lo.CountBy(value,
					func(row db_sql.ListJobsWithProgressRow) bool {
						return row.Status.String == "downloaded" || row.Status.String == "imported"
					})

				imported := lo.CountBy(value,
					func(row db_sql.ListJobsWithProgressRow) bool {
						return row.Status.String == "imported"
					})

				jobsResponse = append(jobsResponse, jobResponse{
					SearchQuery: value[0].SearchQuery.String,
					ID:          key,
					Limit:       int(value[0].Limit.Int64),
					Progesss: jobProgress{
						Imported:   imported,
						Scraped:    scraped,
						Downloaded: downloaded,
					},
				})
			}

			render.JSON(w, r, jobsResponse)
		},
	)
}
