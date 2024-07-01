package api

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/db"
)

func HandleCreateJob(queries *db.Queries, wakeJobs chan<- struct{}) http.HandlerFunc {
	type jobRequest struct {
		SearchQuery string `json:"search_query"`
		Limit       int    `json:"limit"`
		FilterId    string `json:"filter_id"`
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

			res, err := queries.CreateJob(r.Context(), db.CreateJobParams{
				FilterID: sql.NullString{
					String: articleRequest.FilterId,
					Valid:  true,
				},
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

			select {
			case wakeJobs <- struct{}{}:
				slog.Info("Waking up jobs")
			default:
				slog.Info("Jobs already awake")
			}
		},
	)
}

func HandleListJobs(queries *db.Queries) http.HandlerFunc {
	type jobResponse struct {
		ID          string `json:"id"`
		SearchQuery string `json:"search_query"`
		Limit       int    `json:"limit"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			res, err := queries.ListJobsWithVideos(r.Context())
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}

			grouped := lo.GroupBy(res, func(row db.ListJobsWithVideosRow) string {
				return row.ID
			})

			jobsResponse := []jobResponse{}
			for key, value := range grouped {

				if len(value) == 1 && value[0].ID_2.Valid == false {
					jobsResponse = append(jobsResponse, jobResponse{
						SearchQuery: value[0].SearchQuery.String,
						ID:          key,
						Limit:       int(value[0].Limit.Int64),
					})
					continue
				}

				jobsResponse = append(jobsResponse, jobResponse{
					SearchQuery: value[0].SearchQuery.String,
					ID:          key,
					Limit:       int(value[0].Limit.Int64),
				})
			}

			render.JSON(w, r, jobsResponse)
		},
	)
}

func HandleJobDetails(queries *db.Queries) http.HandlerFunc {
	type jobResponse struct {
		ID          string `json:"id"`
		SearchQuery string `json:"search_query"`
		VideoTarget int    `json:"video_target"`
		VideosFound int    `json:"videos_found"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			jobId := chi.URLParam(r, "id")
			if jobId == "" {
				render.Status(r, http.StatusBadRequest)
				render.PlainText(w, r, "No file id provided")
				return
			}

			res, err := queries.GetJob(r.Context(), jobId)
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}

			resp := jobResponse{
				ID:          res.ID,
				SearchQuery: res.SearchQuery.String,
				VideoTarget: int(res.Limit.Int64),
				VideosFound: int(res.VideosFound),
			}

			render.JSON(w, r, resp)
		},
	)
}

func HandleJobVideosFound(queries *db.Queries) http.HandlerFunc {
	type video struct {
		YoutubeId string `json:"youtube_id"`
	}

	type jobResponse struct {
		ID     string  `json:"id"`
		Videos []video `json:"videos"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			jobId := chi.URLParam(r, "id")
			if jobId == "" {
				render.Status(r, http.StatusBadRequest)
				render.PlainText(w, r, "No file id provided")
				return
			}

			job, err := queries.GetOneWithVideos(r.Context(), jobId)
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}

			videos := lo.Map(job, func(row db.GetOneWithVideosRow, index int) video {
				return video{
					YoutubeId: row.VideoYoutubeID.String,
				}
			})

			resp := jobResponse{
				ID:     jobId,
				Videos: videos,
			}

			render.JSON(w, r, resp)
		},
	)
}
