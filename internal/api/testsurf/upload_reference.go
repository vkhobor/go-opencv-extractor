package testsurf

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	u "github.com/vkhobor/go-opencv/internal/api/util"
	"github.com/vkhobor/go-opencv/internal/features/testsurf"
)

type UploadReferenceRequest struct {
	RawBody huma.MultipartFormFiles[struct {
		Video huma.FormFile `form:"reference" contentType:"application/octet-stream" required:"true"`
	}]
}

func HandleUploadReference() u.Handler[UploadVideoRequest, struct{}] {
	return func(ctx context.Context, req *UploadVideoRequest) (*struct{}, error) {
		feat := testsurf.UploadReferenceFeature{}

		err := feat.UploadReference(ctx, req.RawBody.Data().Video.File)
		return &struct{}{}, err
	}
}
