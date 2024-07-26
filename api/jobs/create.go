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

type CreateJobRequest struct {
	Body CreateJob
}

func HandleCreateJob(queries *db.Queries, wakeJobs chan<- struct{}, config config.ServerConfig) u.Handler[CreateJobRequest, CreatedJobResponse] {

	return func(ctx context.Context, wb *CreateJobRequest) (*CreatedJobResponse, error) {

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
