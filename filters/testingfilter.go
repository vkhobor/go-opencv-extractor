package filters

import (
	"iter"

	"github.com/vkhobor/go-opencv/video/videoiter"
)

type TestFilter struct {
}

func (f *TestFilter) SamplingWantFPS() int {
	return 1
}

func (f *TestFilter) FrameFilter(frames iter.Seq2[videoiter.FrameInfo, error]) iter.Seq2[videoiter.FrameInfo, error] {
	return func(yield func(videoiter.FrameInfo, error) bool) {
		firstFrame := true
		for frame, error := range frames {
			if error != nil {
				yield(videoiter.FrameInfo{}, error)
				return
			}

			if firstFrame {
				firstFrame = false
				yield(frame, nil)
				continue
			}
		}
	}
}
