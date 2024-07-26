package jobs

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/samber/lo"
	u "github.com/vkhobor/go-opencv/api/util"
	"github.com/vkhobor/go-opencv/db"
)

type JobVideo struct {
	YoutubeId      string `json:"youtube_id"`
	DownloadStatus string `json:"download_status"`
	ImportStatus   string `json:"import_status"`
}

type JobAndVideos struct {
	ID     string     `json:"id"`
	Videos []JobVideo `json:"videos"`
}

type JobVideosRequest struct {
	ID string `path:"id"`
}

type JobVideosResponse struct {
	Body JobAndVideos
}

func HandleJobVideosFound(queries *db.Queries) u.Handler[JobVideosRequest, JobVideosResponse] {

	return func(ctx context.Context, wpi *JobVideosRequest) (*JobVideosResponse, error) {
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

		return &JobVideosResponse{
			Body: resp,
		}, nil
	}
}
