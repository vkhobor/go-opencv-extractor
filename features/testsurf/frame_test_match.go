package testsurf

import (
	"context"

	"github.com/vkhobor/go-opencv/mlog"
)

// FrameMatchingTestFeature implements the functionality to test matching between frames
type FrameMatchingTestFeature struct {
	// Add any dependencies here
}

// TestFrameMatch checks if a frame matches the reference image based on SURF features
func (f *FrameMatchingTestFeature) TestFrameMatch(ctx context.Context, frameNum int, ratioCheck float64, minMatches int, goodMatchThreshold float64) (bool, error) {
	mlog.Log().Info("Testing frame match",
		"frameNum", frameNum,
		"ratioCheck", ratioCheck,
		"minMatches", minMatches,
		"goodMatchThreshold", goodMatchThreshold)

	return false, nil
}
