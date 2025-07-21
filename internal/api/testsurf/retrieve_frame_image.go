package testsurf

import (
	"context"
	"io"

	"github.com/danielgtaylor/huma/v2"
	"github.com/vkhobor/go-opencv/internal/config"
	"github.com/vkhobor/go-opencv/internal/features/testsurf"
)

type RetrieveFrameImageRequest struct {
	FrameNum int `query:"framenum" required:"true"`
}

func HandleRetrieveFrameImage(config config.DirectoryConfig) func(ctx context.Context, rfir *RetrieveFrameImageRequest) (*huma.StreamResponse, error) {
	return func(ctx context.Context, rfir *RetrieveFrameImageRequest) (*huma.StreamResponse, error) {
		feat := testsurf.RetrieveFrameImageFeature{
			Config: config,
		}

		readCloser, err := feat.GetFrameImage(ctx, rfir.FrameNum)
		if err != nil {
			return nil, huma.Error500InternalServerError("Internal error", err)
		}

		return &huma.StreamResponse{
			Body: func(hctx huma.Context) {
				hctx.SetHeader("Content-Type", "image/jpeg")
				_, err := io.Copy(hctx.BodyWriter(), readCloser)
				if err != nil {
					return
				}
				defer readCloser.Close()
			},
		}, nil
	}
}
