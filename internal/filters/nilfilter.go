package filters

import (
	"iter"

	"github.com/vkhobor/go-opencv/internal/video/videoiter"
)

type NilFilter struct {
}

func (f *NilFilter) SamplingWantFPS() int {
	return 1
}

func (f *NilFilter) FrameFilter(frames iter.Seq2[videoiter.FrameInfo, error]) iter.Seq2[videoiter.FrameInfo, error] {
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
