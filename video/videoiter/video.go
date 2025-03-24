package videoiter

import (
	"errors"
	"math"
	"os"
	"time"

	"github.com/vkhobor/go-opencv/video/metadata"
	"gocv.io/x/gocv"
)

type Video struct {
	path        string
	startFrame  int
	endFrame    int
	originalFPS float64
}

func NewVideo(path string) (Video, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return Video{}, errors.New("Video file does not exist")
	}
	meta, err := metadata.ExtractMetadata(path)
	if err != nil {
		return Video{}, err
	}

	capture, err := gocv.OpenVideoCapture(path)
	if err != nil {
		return Video{}, err
	}
	defer capture.Close()
	maxFrame := int(capture.Get(gocv.VideoCaptureFrameCount))

	iter := Video{
		startFrame:  0,
		endFrame:    maxFrame,
		path:        path,
		originalFPS: meta,
	}

	return iter, nil
}

func (v Video) MaxFrame() int {
	return v.endFrame - v.startFrame
}

func (v Video) FPS() float64 {
	return float64(v.originalFPS)
}

func (v Video) Path() string {
	return v.path
}

func (v Video) GetPercentFrame(frame int) float64 {
	centered := frame - v.startFrame
	return float64(centered) / float64(v.MaxFrame())
}

func (v Video) CurrentTime(frame int) time.Duration {
	return time.Second * time.Duration(math.Ceil(float64(frame)/v.originalFPS))
}

func (v Video) CurrentProgress(frame int) Progress {
	return Progress{
		Done:      frame - v.startFrame,
		Max:       v.MaxFrame(),
		timestamp: time.Now(),
	}
}
