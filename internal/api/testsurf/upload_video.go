package testsurf

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	u "github.com/vkhobor/go-opencv/internal/api/util"
	"github.com/vkhobor/go-opencv/internal/config"
	"github.com/vkhobor/go-opencv/internal/features/testsurf"
)

type UploadVideoRequest struct {
	RawBody huma.MultipartFormFiles[struct {
		Video huma.FormFile `form:"video" contentType:"application/octet-stream" required:"true"`
	}]
}

func HandleUploadVideo(config config.DirectoryConfig) u.Handler[UploadVideoRequest, struct{}] {
	return func(ctx context.Context, req *UploadVideoRequest) (*struct{}, error) {
		feat := testsurf.UploadVideoFeature{
			Config: config,
		}

		err := feat.UploadVideo(ctx, req.RawBody.Data().Video.File)
		return &struct{}{}, err
	}
}
