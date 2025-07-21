package testsurf

import (
	"github.com/vkhobor/go-opencv/internal/config"
	"github.com/vkhobor/go-opencv/internal/video/metadata"
)

// VideoMetadataFeature handles retrieving metadata for test videos
type VideoMetadataFeature struct {
	Config config.DirectoryConfig
}

func (f *VideoMetadataFeature) GetFrameCount() (int, error) {
	return metadata.GetMaxFrames(f.Config.GetTestVideoPath())
}
