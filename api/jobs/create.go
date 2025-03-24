package jobs

import (
	"context"
	"database/sql"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
	u "github.com/vkhobor/go-opencv/api/util"
	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/db"
	"github.com/vkhobor/go-opencv/mlog"
)

type CreateJob struct {
	IsQueryBased  bool `json:"is_query_based"`
	IsSingleVideo bool `json:"is_single_video"`

	SearchQuery string `json:"search_query,omitempty"`
	Limit       int    `json:"limit,omitempty"`
	FilterId    string `json:"filter_id,omitempty"`

	YoutubeId string `json:"youtube_id,omitempty"`
}

type CreatedJob struct {
	Id string `json:"id"`
}

type CreatedJobResponse struct {
	Body     CreatedJob
	Location string `header:"Location"`
}

type CreateJobRequest struct {
	Body CreateJob
}

func HandleCreateJob(queries *db.Queries, wakeJobs chan<- struct{}, config config.ServerConfig) u.Handler[CreateJobRequest, CreatedJobResponse] {

	return func(ctx context.Context, wb *CreateJobRequest) (*CreatedJobResponse, error) {
		if !onlyOneTrue(wb.Body.IsQueryBased, wb.Body.IsSingleVideo) {
			return nil, huma.Error400BadRequest("Exactly one of is_query_based or is_single_video must be true")
		}

		jobs, err := queries.GetJobs(ctx)
		if err != nil {
			return nil, err
		}

		exists := lo.SomeBy(jobs, func(row db.GetJobsRow) bool {
			return (row.SearchQuery.String == wb.Body.SearchQuery ||
				row.YoutubeID.String == wb.Body.YoutubeId) &&
				row.FilterID.String == wb.Body.FilterId
		})
		if exists {
			return nil, huma.Error400BadRequest("Job already exists for filter")
		}

		var res db.Job
		if wb.Body.IsQueryBased {
			res, err = queries.CreateJob(ctx, db.CreateJobParams{
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
		} else if wb.Body.IsSingleVideo {
			res, err = queries.CreateJob(ctx, db.CreateJobParams{
				FilterID: sql.NullString{
					String: wb.Body.FilterId,
					Valid:  true,
				},
				YoutubeID: sql.NullString{
					String: wb.Body.YoutubeId,
					Valid:  true,
				},
				Limit: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				ID: uuid.New().String(),
			})
			if err != nil {
				return nil, err
			}
		}

		select {
		case wakeJobs <- struct{}{}:
			mlog.Log().Info("Waking up jobs")
		default:
			mlog.Log().Info("Jobs already awake")
		}

		return &CreatedJobResponse{
			Body: CreatedJob{
				Id: res.ID,
			},
			Location: config.BaseUrl + "/api/jobs/" + res.ID,
		}, nil
	}
}

func onlyOneTrue(args ...bool) bool {
	count := 0
	for _, arg := range args {
		if arg {
			count++
		}
	}
	return count == 1
}
