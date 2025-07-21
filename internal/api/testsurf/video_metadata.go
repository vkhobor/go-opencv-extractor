package testsurf

import (
	"context"

	u "github.com/vkhobor/go-opencv/internal/api/util"
	"github.com/vkhobor/go-opencv/internal/config"
	"github.com/vkhobor/go-opencv/internal/features/testsurf"
)

type VideoMetadataBody struct {
	MaxFrame int `json:"maxframe"`
}

type VideoMetadataResponse struct {
	Body VideoMetadataBody `json:"body"`
}

func HandleVideoMetadata(config config.DirectoryConfig) u.Handler[struct{}, VideoMetadataResponse] {
	return func(ctx context.Context, _ *struct{}) (*VideoMetadataResponse, error) {
		feat := testsurf.VideoMetadataFeature{
			Config: config,
		}

		frame, err := feat.GetFrameCount()
		if err != nil {
			return nil, err
		}

		return &VideoMetadataResponse{
			Body: VideoMetadataBody{
				MaxFrame: frame,
			},
		}, nil
	}
}
