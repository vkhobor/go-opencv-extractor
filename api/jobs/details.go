package jobs

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	u "github.com/vkhobor/go-opencv/api/util"
	"github.com/vkhobor/go-opencv/db"
)

type JobDetailsRequest struct {
	ID string `path:"id"`
}

type JobDetails struct {
	ID          string `json:"id"`
	SearchQuery string `json:"search_query"`
	VideoTarget int    `json:"video_target"`
	VideosFound int    `json:"videos_found"`
}

type JobDetailsResponse struct {
	Body JobDetails
}

func HandleJobDetails(queries *db.Queries) u.Handler[JobDetailsRequest, JobDetailsResponse] {

	return func(ctx context.Context, wpi *JobDetailsRequest) (*JobDetailsResponse, error) {
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
		return &JobDetailsResponse{
			Body: resp,
		}, nil
	}
}
