package filters

import (
	"iter"
	"time"

	mse "github.com/vkhobor/go-opencv/image/mean_squared_error"
	"github.com/vkhobor/go-opencv/image/surf"
	"github.com/vkhobor/go-opencv/video/videoiter"
)

type SURFVideoFilter struct {
	surfMatcher *surf.SURFImageMatcher
	// Good default is 0.2
	MSEThreshold float64
}

func NewSURFVideoFilter(surfMatcher *surf.SURFImageMatcher, MSEThreshold float64) *SURFVideoFilter {
	return &SURFVideoFilter{
		surfMatcher:  surfMatcher,
		MSEThreshold: MSEThreshold,
	}
}

func (f *SURFVideoFilter) SamplingWantFPS() int {
	return 1
}

func (f *SURFVideoFilter) FrameFilter(frames iter.Seq2[videoiter.FrameInfo, error]) iter.Seq2[videoiter.FrameInfo, error] {
	previousFrame := videoiter.FrameInfo{}
	var firstFrame bool = true

	return func(yield func(videoiter.FrameInfo, error) bool) {
		for frame, error := range frames {
			if firstFrame {
				previousFrame = frame.Clone()
				firstFrame = false
				continue
			}

			if error != nil {
				yield(videoiter.FrameInfo{}, error)
				// Abort iteration on error
				return
			}

			// If time distance is large enough, we expect frames to be different
			// so worth checking if they are what we need
			if distanceIsMoreThanDuration(previousFrame, frame, time.Minute*2) {
				if f.surfMatcher.IsImageMatch(&frame.Frame) {
					previousFrame = frame.Clone()
					yield(frame, nil)
					continue
				}
			}

			// If time distance is small, we expect to get the exact same frame,
			// so it is worth skipping if they are super similar with the previous frame
			diff, err := mse.GetMeanSquaredError(&previousFrame.Frame, &frame.Frame)
			if err != nil {
				yield(videoiter.FrameInfo{}, err)
				// Abort iteration on error
				return
			}
			if diff < f.MSEThreshold {
				previousFrame = frame.Clone()
				continue
			}

			// We got here because the frames were close in time,
			// but probably the camera cut so worth checking that we need it or not
			if f.surfMatcher.IsImageMatch(&frame.Frame) {
				previousFrame = frame.Clone()
				yield(frame, nil)
				continue
			}
		}
	}
}

func distanceIsMoreThanDuration(one, two videoiter.FrameInfo, duration time.Duration) bool {
	dist := one.TimeFromStart - two.TimeFromStart
	return dist.Abs() >= duration
}
