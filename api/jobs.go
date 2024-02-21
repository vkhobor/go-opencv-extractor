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
		Imported         int      `json:"imported"`
		Scraped          int      `json:"scraped"`
		Downloaded       int      `json:"downloaded"`
		VideoIds         []string `json:"video_ids"`
		NumberOfPictures int      `json:"number_of_pictures"`
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
							Imported:         0,
							Scraped:          0,
							Downloaded:       0,
							VideoIds:         []string{},
							NumberOfPictures: 0,
						},
					})
					continue
				}

				downloaded := lo.Filter(value,
					func(row db_sql.ListJobsWithProgressRow, index int) bool {
						return row.BlobStorageID.Valid
					})
				downloaded = lo.UniqBy(downloaded, func(item db_sql.ListJobsWithProgressRow) string {
					return item.BlobStorageID.String
				})

				imported := lo.Filter(value,
					func(row db_sql.ListJobsWithProgressRow, index int) bool {
						return row.Status.String == "imported"
					})
				imported = lo.UniqBy(imported, func(item db_sql.ListJobsWithProgressRow) string {
					return item.ID_2.String
				})

				pictures := lo.Filter(value, func(item db_sql.ListJobsWithProgressRow, i int) bool {
					return item.ID_3.Valid
				})
				pictures = lo.UniqBy(pictures, func(item db_sql.ListJobsWithProgressRow) string {
					return item.ID_3.String
				})

				allVideoIds := lo.Map(
					lo.Filter(value, func(item db_sql.ListJobsWithProgressRow, index int) bool {
						return item.ID_2.Valid
					}), func(item db_sql.ListJobsWithProgressRow, i int) string {
						return item.ID_2.String
					})

				uniqueVideoIds := lo.Uniq(allVideoIds)

				jobsResponse = append(jobsResponse, jobResponse{
					SearchQuery: value[0].SearchQuery.String,
					ID:          key,
					Limit:       int(value[0].Limit.Int64),
					Progesss: jobProgress{
						Imported:         len(imported),
						Scraped:          len(uniqueVideoIds),
						Downloaded:       len(downloaded),
						NumberOfPictures: len(pictures),
						VideoIds:         uniqueVideoIds,
					},
				})
			}

			render.JSON(w, r, jobsResponse)
		},
	)
}
