package metadata

import (
	"errors"
	"os"

	"gocv.io/x/gocv"
)

func GetMaxFrames(path string) (int, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return 0, errors.New("Video file does not exist")
	}

	capture, err := gocv.OpenVideoCapture(path)
	if err != nil {
		return 0, err
	}
	defer capture.Close()
	maxFrame := int(capture.Get(gocv.VideoCaptureFrameCount))

	return maxFrame, nil
}
