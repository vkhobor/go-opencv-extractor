package api

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/db_sql"
)

func HandleCreateJob(queries *db_sql.Queries, wakeJobs chan<- struct{}) http.HandlerFunc {
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

			select {
			case wakeJobs <- struct{}{}:
				slog.Info("Waking up jobs")
			default:
				slog.Info("Jobs already awake")
			}
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

func HandleJobDetails(queries *db_sql.Queries) http.HandlerFunc {
	type jobResponse struct {
		ID            string `json:"id"`
		SearchQuery   string `json:"search_query"`
		VideoTarget   int    `json:"video_target"`
		PicturesFound int    `json:"pictures_found"`
		VideosFound   int    `json:"videos_found"`
		VideosInError int    `json:"videos_in_error"`
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
				ID:            res.ID,
				SearchQuery:   res.SearchQuery.String,
				VideoTarget:   int(res.Limit.Int64),
				PicturesFound: int(res.PicturesFound),
				VideosFound:   int(res.VideosFound),
				VideosInError: int(res.VideosInError),
			}

			render.JSON(w, r, resp)
		},
	)
}

func HandleJobVideosFound(queries *db_sql.Queries) http.HandlerFunc {
	type video struct {
		YoutubeId     string `json:"youtube_id"`
		Status        string `json:"status"`
		PicturesFound int    `json:"pictures_found"`
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

			res, err := queries.GetOneWithVideos(r.Context(), jobId)
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}

			videos := lo.Map(res, func(row db_sql.GetOneWithVideosRow, index int) video {
				return video{
					YoutubeId:     row.VideoYoutubeID.String,
					Status:        row.VideoStatus.String,
					PicturesFound: int(row.PicturesFound),
				}
			})

			resp := jobResponse{
				ID:     res[0].ID,
				Videos: videos,
			}

			render.JSON(w, r, resp)
		},
	)
}

func HandleJobProgress(queries *db_sql.Queries) http.HandlerFunc {
	type jobResponse struct {
		ID               string `json:"id"`
		Imported         int    `json:"imported"`
		Scraped          int    `json:"scraped"`
		Downloaded       int    `json:"downloaded"`
		NumberOfPictures int    `json:"number_of_pictures"`
		Limit            int    `json:"limit"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			jobId := chi.URLParam(r, "id")
			if jobId == "" {
				render.Status(r, http.StatusBadRequest)
				render.PlainText(w, r, "No file id provided")
				return
			}

			res, err := queries.GetJobWithProgress(r.Context(), jobId)
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}

			resp := jobResponse{
				ID:               res.ID,
				Imported:         int(res.VideosImported),
				Scraped:          int(res.VideosScraped) + int(res.VideosDownloaded) + int(res.VideosImported),
				Downloaded:       int(res.VideosDownloaded) + int(res.VideosImported),
				NumberOfPictures: int(res.PicturesFound),
				Limit:            int(res.Limit.Int64),
			}

			render.JSON(w, r, resp)
		},
	)
}
