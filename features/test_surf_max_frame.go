package features

import (
	"context"

	"github.com/vkhobor/go-opencv/db"
)

type TestSurfMaxFrameFeature struct {
	Queries *db.Queries
}

func (f *TestSurfMaxFrameFeature) GetMaxFrame(ctx context.Context) (int, error) {
	// Logic to return the maximum frame
	return 250000, nil
}
