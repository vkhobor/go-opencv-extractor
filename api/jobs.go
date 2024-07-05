package api

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/db"
)

type CreateJob struct {
	SearchQuery string `json:"search_query"`
	Limit       int    `json:"limit"`
	FilterId    string `json:"filter_id"`
}

type CreatedJob struct {
	Id string `json:"id"`
}

type CreatedJobResponse struct {
	Body     CreatedJob
	Location string `header:"Location"`
}

func HandleCreateJob(queries *db.Queries, wakeJobs chan<- struct{}, config config.ProgramConfig) Handler[WithBody[CreateJob], CreatedJobResponse] {

	return func(ctx context.Context, wb *WithBody[CreateJob]) (*CreatedJobResponse, error) {

		jobs, err := queries.GetJobs(ctx)
		if err != nil {
			return nil, err
		}

		exists := lo.SomeBy(jobs, func(row db.GetJobsRow) bool {
			return row.SearchQuery.String == wb.Body.SearchQuery && row.FilterID.String == wb.Body.FilterId
		})
		if exists {
			return nil, huma.Error400BadRequest("Job already exists for query")
		}

		res, err := queries.CreateJob(ctx, db.CreateJobParams{
			FilterID: sql.NullString{
				String: wb.Body.FilterId,
				Valid:  true,
			},
			SearchQuery: sql.NullString{
				String: wb.Body.SearchQuery,
				Valid:  true,
			},
			Limit: sql.NullInt64{
				Int64: int64(wb.Body.Limit),
				Valid: true,
			},
			ID: uuid.New().String(),
		})

		if err != nil {
			return nil, err
		}

		select {
		case wakeJobs <- struct{}{}:
			slog.Info("Waking up jobs")
		default:
			slog.Info("Jobs already awake")
		}

		return &CreatedJobResponse{
			Body: CreatedJob{
				Id: res.ID,
			},
			Location: config.BaseUrl + "/api/jobs/" + res.ID,
		}, nil
	}
}

func HandleRestartJobPipeline(wakeJobs chan<- struct{}) Handler[Empty, Empty] {

	return func(ctx context.Context, e *Empty) (*Empty, error) {
		go func() {
			wakeJobs <- struct{}{}
		}()
		return &Empty{}, nil
	}
}

type ListJobResponse struct {
	ID          string `json:"id"`
	SearchQuery string `json:"search_query"`
	Limit       int    `json:"limit"`
}

func HandleListJobs(queries *db.Queries) Handler[Empty, WithBody[[]ListJobResponse]] {

	return func(ctx context.Context, e *Empty) (*WithBody[[]ListJobResponse], error) {
		res, err := queries.ListJobsWithVideos(ctx)
		if err != nil {
			return nil, err
		}

		grouped := lo.GroupBy(res, func(row db.ListJobsWithVideosRow) string {
			return row.ID
		})

		jobsResponse := []ListJobResponse{}
		for key, value := range grouped {

			if len(value) == 1 && value[0].ID_2.Valid == false {
				jobsResponse = append(jobsResponse, ListJobResponse{
					SearchQuery: value[0].SearchQuery.String,
					ID:          key,
					Limit:       int(value[0].Limit.Int64),
				})
				continue
			}

			jobsResponse = append(jobsResponse, ListJobResponse{
				SearchQuery: value[0].SearchQuery.String,
				ID:          key,
				Limit:       int(value[0].Limit.Int64),
			})
		}
		return &WithBody[[]ListJobResponse]{
			Body: jobsResponse,
		}, nil
	}
}

type JobDetails struct {
	ID          string `json:"id"`
	SearchQuery string `json:"search_query"`
	VideoTarget int    `json:"video_target"`
	VideosFound int    `json:"videos_found"`
}

func HandleJobDetails(queries *db.Queries) Handler[WithPathId, WithBody[JobDetails]] {

	return func(ctx context.Context, wpi *WithPathId) (*WithBody[JobDetails], error) {
		if wpi.ID == "" {
			return nil, huma.Error400BadRequest("id not found")
		}

		res, err := queries.GetJob(ctx, wpi.ID)
		if err != nil {
			return nil, err
		}

		resp := JobDetails{
			ID:          res.ID,
			SearchQuery: res.SearchQuery.String,
			VideoTarget: int(res.Limit.Int64),
			VideosFound: int(res.VideosFound),
		}
		return &WithBody[JobDetails]{
			Body: resp,
		}, nil
	}
}

type JobVideo struct {
	YoutubeId      string `json:"youtube_id"`
	DownloadStatus string `json:"download_status"`
	ImportStatus   string `json:"import_status"`
}

type JobAndVideos struct {
	ID     string     `json:"id"`
	Videos []JobVideo `json:"videos"`
}

func HandleJobVideosFound(queries *db.Queries) Handler[WithPathId, WithBody[JobAndVideos]] {

	return func(ctx context.Context, wpi *WithPathId) (*WithBody[JobAndVideos], error) {
		if wpi.ID == "" {
			return nil, huma.Error400BadRequest("id not found")
		}

		job, err := queries.GetVideosForJob(ctx, wpi.ID)
		if err != nil {
			return nil, err
		}

		videos := lo.Map(job, func(row db.GetVideosForJobRow, index int) JobVideo {
			downloadStatus := "not started or progressing"
			if row.DownloadAttemptsSuccess > 0 {
				downloadStatus = "success"
			} else if row.DownloadAttemptsError > 0 {
				downloadStatus = "failed"
			}

			importStatus := "not started or progressing"
			if row.ImportAttemptsSuccess > 0 {
				importStatus = "success"
			} else if row.ImportAttemptsError > 0 {
				importStatus = "failed"
			}

			return JobVideo{
				YoutubeId:      row.VideoYoutubeID,
				DownloadStatus: downloadStatus,
				ImportStatus:   importStatus,
			}
		})

		resp := JobAndVideos{
			ID:     wpi.ID,
			Videos: videos,
		}

		return &WithBody[JobAndVideos]{
			Body: resp,
		}, nil
	}
}
