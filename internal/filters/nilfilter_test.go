package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vkhobor/go-opencv/internal/video/videoiter"
)

func TestNilFilter(t *testing.T) {
	filter := &NilFilter{}

	video, err := videoiter.NewVideo("../../samples/7WAkEo9i6ts.mp4")
	assert.NoError(t, err)

	fpsWant := filter.SamplingWantFPS()
	frames := videoiter.AllSampledFrames(video, fpsWant)
	wantFrames := filter.FrameFilter(frames)

	found := false
	for frame, err := range wantFrames {
		assert.NoError(t, err)
		if frame.Frame.Empty() {
			continue
		}
		found = true
		break
	}

	assert.True(t, found, "NilFilter should find at least one frame")
}
