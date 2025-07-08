package testsurf

import (
	"context"
	"errors"

	"github.com/vkhobor/go-opencv/image/surf"
	"github.com/vkhobor/go-opencv/mlog"
	"gocv.io/x/gocv"
)

// FrameMatchingTestFeature implements the functionality to test matching between frames
type FrameMatchingTestFeature struct {
	// Add any dependencies here
}

// TestFrameMatch checks if a frame matches the reference image based on SURF features
func (f *FrameMatchingTestFeature) TestFrameMatch(ctx context.Context, frameNum int, ratioCheck float64, minMatches int, goodMatchThreshold float64) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	mlog.Log().Info("Testing frame match",
		"frameNum", frameNum,
		"ratioCheck", ratioCheck,
		"minMatches", minMatches,
		"goodMatchThreshold", goodMatchThreshold)

	if err := ctx.Err(); err != nil {
		return false, err
	}

	if frameNum < 0 {
		return false, errors.New("Frame number must be non-negative")
	}

	options := []surf.SURFImageMatcherOption{
		surf.WithMinMatches(int(minMatches)),
		surf.WithMinThreshold(goodMatchThreshold),
		surf.WithRatioThreshold(ratioCheck),
	}

	if err := ctx.Err(); err != nil {
		return false, err
	}

	if cachedReferenceImage.Empty() {
		return false, errors.New("Reference image is empty")
	}

	matcher, err := surf.NewSURFImageMatcherFromMats([]gocv.Mat{*cachedReferenceImage}, options...)
	if err != nil {
		return false, err
	}
	defer matcher.Close()

	if err := ctx.Err(); err != nil {
		return false, err
	}

	frame, err := cachedTestVideoExtractor.GetFrameAsMat(ctx, frameNum)
	if err != nil {
		return false, err
	}
	defer frame.Close()

	if err := ctx.Err(); err != nil {
		return false, err
	}

	matched := matcher.IsImageMatch(&frame)
	return matched, nil
}
