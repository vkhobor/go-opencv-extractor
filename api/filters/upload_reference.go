package filters

import (
	"context"
	"database/sql"

	"github.com/danielgtaylor/huma/v2"
	u "github.com/vkhobor/go-opencv/api/util"
	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/features"
)

const (
	megabyte = 1 << 20 // 1 megabyte = 2^20 bytes
)

type ReferenceUploadRequest struct {
	RawBody huma.MultipartFormFiles[struct {
		File                       huma.FormFile `form:"file" required:"true"`
		RatioTestThreshold         float64       `form:"ratioTestThreshold" required:"true"`
		MinThresholdForSURFMatches float64       `form:"minThresholdForSURFMatches" required:"true"`
		MinSURFMatches             int64         `form:"minSURFMatches" required:"true"`
		MseSkip                    float64       `form:"mseSkip" required:"true"`
	}]
}

type ReferenceUploadResponse struct {
	Body ReferenceUploadResponseBody `json:"body"`
}

type ReferenceUploadResponseBody struct {
	Status string `json:"status"`
}

// TODO migrate to fully dynamic filters
var defaultFilterId = "1fed33d4-0ea3-4b84-909c-261e4b2a3d43"

func HandleReferenceUpload(queries *sql.DB, config config.DirectoryConfig) u.Handler[ReferenceUploadRequest, ReferenceUploadResponse] {
	dbAdapter := u.NewDbAdapter(queries)

	feature := &features.ReferenceUploadFeature{
		Querier: dbAdapter.Querier,
		SqlDB:   dbAdapter.TxEr,
		Config:  config,
	}

	return func(ctx context.Context, req *ReferenceUploadRequest) (*ReferenceUploadResponse, error) {
		data := req.RawBody.Data()

		err := feature.UploadReference(ctx, data.File.File, data.File.Filename, features.ReferenceConfig{
			RatioTestThreshold:         data.RatioTestThreshold,
			MinThresholdForSURFMatches: data.MinThresholdForSURFMatches,
			MinSURFMatches:             data.MinSURFMatches,
			MseSkip:                    data.MseSkip,
		})
		if err != nil {
			return nil, err
		}

		return &ReferenceUploadResponse{Body: ReferenceUploadResponseBody{
			Status: "y",
		}}, nil
	}
}
