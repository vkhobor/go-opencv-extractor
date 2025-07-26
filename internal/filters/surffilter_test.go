package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vkhobor/go-opencv/internal/image/surf"
	"github.com/vkhobor/go-opencv/internal/video/videoiter"
)

func TestSURFVideoFilter(t *testing.T) {
	matcher, err := surf.NewSURFImageMatcher([]string{"../../samples/vertic.png"})
	assert.NoError(t, err)
	defer matcher.Close()

	filter := NewSURFVideoFilter(matcher, 0.2)

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

	assert.True(t, found, "SURFVideoFilter should find at least one matching frame")
}
