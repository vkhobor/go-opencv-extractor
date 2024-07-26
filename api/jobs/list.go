package jobs

import (
	"context"

	"github.com/samber/lo"
	u "github.com/vkhobor/go-opencv/api/util"
	"github.com/vkhobor/go-opencv/db"
)

type ListJobBody struct {
	ID          string `json:"id"`
	SearchQuery string `json:"search_query"`
	Limit       int    `json:"limit"`
}

type ListJobResponse struct {
	Body []ListJobBody
}

func HandleListJobs(queries *db.Queries) u.Handler[struct{}, ListJobResponse] {

	return func(ctx context.Context, e *struct{}) (*ListJobResponse, error) {
		res, err := queries.ListJobsWithVideos(ctx)
		if err != nil {
			return nil, err
		}

		grouped := lo.GroupBy(res, func(row db.ListJobsWithVideosRow) string {
			return row.ID
		})

		jobsResponse := []ListJobBody{}
		for key, value := range grouped {

			if len(value) == 1 && value[0].ID_2.Valid == false {
				jobsResponse = append(jobsResponse, ListJobBody{
					SearchQuery: value[0].SearchQuery.String,
					ID:          key,
					Limit:       int(value[0].Limit.Int64),
				})
				continue
			}

			jobsResponse = append(jobsResponse, ListJobBody{
				SearchQuery: value[0].SearchQuery.String,
				ID:          key,
				Limit:       int(value[0].Limit.Int64),
			})
		}
		return &ListJobResponse{
			Body: jobsResponse,
		}, nil
	}
}
