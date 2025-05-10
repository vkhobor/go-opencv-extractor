package testsurf

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	u "github.com/vkhobor/go-opencv/api/util"
)

type UploadVideoRequest struct {
	RawBody huma.MultipartFormFiles[struct {
		Video huma.FormFile `form:"video" contentType:"video/*" required:"true"`
	}]
}

func HandleUploadVideo() u.Handler[UploadVideoRequest, struct{}] {
	return func(ctx context.Context, req *UploadVideoRequest) (*struct{}, error) {
		return &struct{}{}, nil
	}
}
