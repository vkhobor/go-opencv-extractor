package features

import (
	"context"

	"github.com/vkhobor/go-opencv/db"
)

type TestSurfFrameFeature struct {
	Queries *db.Queries
}

func (f *TestSurfFrameFeature) GetFrame(ctx context.Context, frameNum int) (string, error) {
	// Logic to serve an image for a specific frame number
	return "Image data for frame", nil
}
