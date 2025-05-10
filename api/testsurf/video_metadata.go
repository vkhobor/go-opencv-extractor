package testsurf

import (
	"context"

	u "github.com/vkhobor/go-opencv/api/util"
)

type VideoMetadataBody struct {
	MaxFrame int `json:"maxframe"`
}

type VideoMetadataResponse struct {
	Body VideoMetadataBody `json:"body"`
}

func HandleVideoMetadata() u.Handler[struct{}, VideoMetadataResponse] {
	return func(ctx context.Context, _ *struct{}) (*VideoMetadataResponse, error) {
		return &VideoMetadataResponse{
			Body: VideoMetadataBody{
				MaxFrame: 0,
			},
		}, nil
	}
}
