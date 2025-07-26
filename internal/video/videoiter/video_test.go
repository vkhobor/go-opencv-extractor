package videoiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewVideo(t *testing.T) {
	t.Run("NonExistentFile", func(t *testing.T) {
		_, err := NewVideo("non_existent_file.mp4")
		assert.Error(t, err, "Expected error for non-existent file")
	})

	t.Run("ValidFile", func(t *testing.T) {
		video, err := NewVideo("../../../samples/7WAkEo9i6ts.mp4")
		assert.NoError(t, err, "Expected no error for valid file")
		assert.Equal(t, "../../../samples/7WAkEo9i6ts.mp4", video.Path(), "Video path mismatch")
	})
}

func TestMaxFrame(t *testing.T) {
	video := Video{
		path:        "mock_video.mp4",
		startFrame:  0,
		endFrame:    100,
		originalFPS: 30.0,
	}

	t.Run("MaxFrame", func(t *testing.T) {
		assert.Equal(t, 100, video.MaxFrame(), "MaxFrame calculation is incorrect")
	})
}

func TestFPS(t *testing.T) {
	video := Video{
		path:        "mock_video.mp4",
		startFrame:  0,
		endFrame:    100,
		originalFPS: 30.0,
	}

	t.Run("FPS", func(t *testing.T) {
		assert.Equal(t, 30.0, video.FPS(), "FPS value is incorrect")
	})
}

func TestPath(t *testing.T) {
	video := Video{
		path:        "mock_video.mp4",
		startFrame:  0,
		endFrame:    100,
		originalFPS: 30.0,
	}

	t.Run("Path", func(t *testing.T) {
		assert.Equal(t, "mock_video.mp4", video.Path(), "Path value is incorrect")
	})
}

func TestGetPercentFrame(t *testing.T) {
	tests := []struct {
		name        string
		video       Video
		frame       int
		wantPercent float64
	}{
		{"MiddleFrame", Video{path: "mock_video.mp4", startFrame: 0, endFrame: 100, originalFPS: 30.0}, 50, 0.5},
		{"StartFrame", Video{path: "mock_video.mp4", startFrame: 0, endFrame: 100, originalFPS: 30.0}, 0, 0.0},
		{"EndFrame", Video{path: "mock_video.mp4", startFrame: 0, endFrame: 100, originalFPS: 30.0}, 100, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantPercent, tt.video.GetPercentFrame(tt.frame), "GetPercentFrame calculation is incorrect")
		})
	}
}

func TestCurrentTime(t *testing.T) {
	tests := []struct {
		name         string
		video        Video
		frame        int
		wantDuration time.Duration
	}{
		{"MiddleFrame", Video{path: "mock_video.mp4", startFrame: 0, endFrame: 100, originalFPS: 30.0}, 60, time.Second * 2},
		{"StartFrame", Video{path: "mock_video.mp4", startFrame: 0, endFrame: 100, originalFPS: 30.0}, 0, time.Second * 0},
		{"EndFrame", Video{path: "mock_video.mp4", startFrame: 0, endFrame: 100, originalFPS: 30.0}, 90, time.Second * 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantDuration, tt.video.CurrentTime(tt.frame), "CurrentTime calculation is incorrect")
		})
	}
}

func TestCurrentProgress(t *testing.T) {
	tests := []struct {
		name     string
		video    Video
		frame    int
		wantDone int
		wantMax  int
	}{
		{"MiddleFrame", Video{path: "mock_video.mp4", startFrame: 0, endFrame: 100, originalFPS: 30.0}, 50, 50, 100},
		{"StartFrame", Video{path: "mock_video.mp4", startFrame: 0, endFrame: 100, originalFPS: 30.0}, 0, 0, 100},
		{"EndFrame", Video{path: "mock_video.mp4", startFrame: 0, endFrame: 100, originalFPS: 30.0}, 100, 100, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			progress := tt.video.CurrentProgress(tt.frame)
			assert.Equal(t, tt.wantDone, progress.Done, "Progress Done value is incorrect")
			assert.Equal(t, tt.wantMax, progress.Max, "Progress Max value is incorrect")
		})
	}
}
