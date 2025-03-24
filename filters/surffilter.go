package filters

import (
	"iter"
	"time"

	mse "github.com/vkhobor/go-opencv/image/mean_squared_error"
	"github.com/vkhobor/go-opencv/image/surf"
	"github.com/vkhobor/go-opencv/video/videoiter"
	"gocv.io/x/gocv"
)

type SURFVideoFilter struct {
	surfMatcher *surf.SURFImageMatcher
}

func NewSURFVideoFilter(surfMatcher *surf.SURFImageMatcher) *SURFVideoFilter {
	return &SURFVideoFilter{
		surfMatcher: surfMatcher,
	}
}

func (f *SURFVideoFilter) SamplingWantFPS() int {
	return 1
}

func (f *SURFVideoFilter) FrameFilter(frames iter.Seq2[videoiter.FrameInfo, error]) iter.Seq2[videoiter.FrameInfo, error] {
	previousFrame := videoiter.FrameInfo{
		Frame: gocv.NewMat(),
	}
	var firstFrame bool = true

	return func(yield func(videoiter.FrameInfo, error) bool) {
		for frame, error := range frames {
			if firstFrame {
				previousFrame.Frame = frame.Frame.Clone()
				previousFrame.FrameNum = frame.FrameNum
				previousFrame.TimeFromStart = frame.TimeFromStart
				firstFrame = false
				continue
			} else {
				previousFrame = frame
			}

			if error != nil {
				yield(videoiter.FrameInfo{}, error)
				return
			}

			if distanceIsLessThanDuration(previousFrame, frame, time.Minute*2) {
				diff, err := mse.GetMeanSquaredError(&previousFrame.Frame, &frame.Frame)
				if err != nil {
					yield(videoiter.FrameInfo{}, err)
					return
				}

				if diff < 0.2 {
					continue
				}
			}

			if !f.surfMatcher.IsImageMatch(&frame.Frame) {
				yield(videoiter.FrameInfo{}, nil)
				return
			}

			if !yield(frame, nil) {
				return
			}
		}
	}
}

func distanceIsLessThanDuration(one, two videoiter.FrameInfo, duration time.Duration) bool {
	dist := one.TimeFromStart - two.TimeFromStart
	return dist.Abs() < duration
}
