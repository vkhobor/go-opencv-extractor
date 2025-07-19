package jobs

import (
	"context"
	"database/sql"
	"path/filepath"

	"github.com/danielgtaylor/huma/v2"
	u "github.com/vkhobor/go-opencv/api/util"
	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/features"
)

type ImportJob struct {
	FilterId string `json:"filter_id" form:"filter_id"`
}

type ImportJobResponse struct {
	Body     CreatedImportJob `json:"body"`
	Location string           `header:"Location"`
}

type CreatedImportJob struct {
	Id string `json:"id"`
}

type ImportJobRequest struct {
	RawBody huma.MultipartFormFiles[struct {
		Video    huma.FormFile `form:"file" contentType:"text/plain" required:"true"`
		FilterId string        `form:"filter_id" required:"true"`
		Name     string        `form:"name"`
	}]
}

func HandleImportJob(q *sql.DB, config config.ServerConfig, wakeJobs chan<- struct{}) u.Handler[ImportJobRequest, ImportJobResponse] {
	return func(ctx context.Context, wb *ImportJobRequest) (*ImportJobResponse, error) {
		conf, err := config.GetDirectoryConfig()
		if err != nil {
			return nil, err
		}
		adapt := u.NewDbAdapter(q)
		uploadFeature := features.UploadVideoFeature{
			DbSql:    adapt.TxEr,
			Querier:  adapt.Querier,
			Config:   conf,
			WakeJobs: wakeJobs,
		}
		data := wb.RawBody.Data()
		path, err := uploadFeature.DownloadVideo(ctx, data.Video.File, data.FilterId, data.Name)
		if err != nil {
			return nil, err
		}

		return &ImportJobResponse{
			Body: CreatedImportJob{
				Id: filepath.Base(path),
			},
			Location: path,
		}, nil
	}
}
