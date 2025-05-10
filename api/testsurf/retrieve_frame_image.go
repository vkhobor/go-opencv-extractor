package testsurf

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

type RetrieveFrameImageRequest struct {
	FrameNum int `query:"framenum" required:"true"`
}

func HandleRetrieveFrameImage() func(ctx context.Context, rfir *RetrieveFrameImageRequest) (*huma.StreamResponse, error) {
	return func(ctx context.Context, rfir *RetrieveFrameImageRequest) (*huma.StreamResponse, error) {
		return &huma.StreamResponse{
			Body: func(hctx huma.Context) {
				hctx.SetHeader("Content-Type", "image/jpeg")
				writer := hctx.BodyWriter()

				// Example: Write a placeholder image or stream frame data
				writer.Write([]byte("Streaming frame image data..."))
			},
		}, nil
	}
}
