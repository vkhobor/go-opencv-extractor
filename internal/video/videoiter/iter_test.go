package videoiter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllSampledFrames(t *testing.T) {
	video := Video{
		path:        "../../../samples/7WAkEo9i6ts.mp4",
		startFrame:  0,
		endFrame:    10,
		originalFPS: 30.0,
	}

	var frames []FrameInfo

	seq := AllSampledFrames(video, 30)
	for frame, _ := range seq {
		frames = append(frames, frame.Clone())
	}

	assert.NotEmpty(t, frames, "frames should not be empty")
	assert.Equal(t, 11, len(frames), "number of frames should be equal of all frames 0-10 is 11 frame")
}
