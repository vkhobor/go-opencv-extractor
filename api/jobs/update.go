package jobs

import (
	"context"
	"database/sql"

	"github.com/danielgtaylor/huma/v2"
	u "github.com/vkhobor/go-opencv/api/util"
	"github.com/vkhobor/go-opencv/db"
	"github.com/vkhobor/go-opencv/mlog"
)

type UpdateJobLimitRequest struct {
	Body struct {
		Limit int `json:"limit"`
	}
	ID string `path:"id"`
}

func HandleUpdateJobLimit(queries *db.Queries, wakeJobs chan<- struct{}) u.Handler[UpdateJobLimitRequest, struct{}] {

	return func(ctx context.Context, wb *UpdateJobLimitRequest) (*struct{}, error) {

		job, err := queries.GetJob(ctx, wb.ID)
		if err != nil {
			return nil, err
		}

		if int64(wb.Body.Limit) <= job.Limit.Int64 {
			return nil, huma.Error400BadRequest("Limit is already set to this value or larger")
		}

		if job.YoutubeID.Valid {
			return nil, huma.Error400BadRequest("Limit is fixed to 1 for simple jobs")
		}

		err = queries.UpdateJobLimit(ctx, db.UpdateJobLimitParams{
			Limit: sql.NullInt64{
				Int64: int64(wb.Body.Limit),
				Valid: true,
			},
			ID: wb.ID,
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

		return &struct{}{}, nil
	}
}

func HandleRestartJobPipeline(wakeJobs chan<- struct{}) u.Handler[struct{}, struct{}] {
	return func(ctx context.Context, e *struct{}) (*struct{}, error) {
		go func() {
			wakeJobs <- struct{}{}
		}()
		return &struct{}{}, nil
	}
}
