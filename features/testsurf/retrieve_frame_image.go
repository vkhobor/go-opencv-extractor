package testsurf

import (
	"context"
	"io"

	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/mlog"
	"github.com/vkhobor/go-opencv/video"
)

var cachedTestVideoExtractor video.FrameExtractor

type RetrieveFrameImageFeature struct {
	Config config.DirectoryConfig
}

func (f *RetrieveFrameImageFeature) GetFrameImage(ctx context.Context, frameNum int) (io.ReadCloser, error) {
	mlog.Log().Info("Retrieving frame image", "frameNum", frameNum)

	frame, err := cachedTestVideoExtractor.GetFrame(frameNum)
	if err != nil {
		return nil, err
	}

	return frame, nil
}
