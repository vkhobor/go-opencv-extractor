package testsurf

import (
	"context"
	"io"

	"github.com/vkhobor/go-opencv/internal/config"
	"github.com/vkhobor/go-opencv/internal/mlog"
)

type RetrieveFrameImageFeature struct {
	Config config.DirectoryConfig
}

func (f *RetrieveFrameImageFeature) GetFrameImage(ctx context.Context, frameNum int) (io.ReadCloser, error) {
	mlog.Log().Info("Retrieving frame image", "frameNum", frameNum)

	frame, err := cachedTestVideoExtractor.GetFrameAsJpeg(ctx, frameNum)
	if err != nil {
		return nil, err
	}

	return frame, nil
}
