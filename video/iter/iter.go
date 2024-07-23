package iter

import (
	"errors"
	"math"
	"os"
	"time"

	"github.com/vkhobor/go-opencv/video/metadata"
	"gocv.io/x/gocv"
)

type VideoIterator struct {
	path        string
	sampleRate  int
	originalFPS float64
	maxFrame    int
}

type FrameInfo struct {
	Frame         *gocv.Mat
	FrameNum      int
	TimeFromStart time.Duration
}

type FrameInfoWithPrevious struct {
	Current  FrameInfo
	Previous FrameInfo
}

func NewVideoIterator(path string, sampleRate int) (VideoIterator, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return VideoIterator{}, errors.New("Video file does not exist")
	}
	meta, err := metadata.ExtractMetadata(path)
	if err != nil {
		return VideoIterator{}, err
	}

	capture, err := gocv.OpenVideoCapture(path)
	if err != nil {
		return VideoIterator{}, err
	}
	defer capture.Close()
	maxFrame := int(capture.Get(gocv.VideoCaptureFrameCount))

	iter := VideoIterator{
		path:        path,
		sampleRate:  sampleRate,
		originalFPS: meta,
		maxFrame:    maxFrame,
	}

	return iter, nil
}

func (v VideoIterator) MaxFrame() int {
	return v.maxFrame
}

func (v VideoIterator) moduloToAchieveTargetFps() int {
	return int(math.Ceil(float64(v.originalFPS) / float64(v.sampleRate)))
}

func (v VideoIterator) currentTime(frame int) time.Duration {
	return time.Second * time.Duration(math.Ceil(float64(frame)/v.originalFPS))
}

func (v VideoIterator) Iterate(yield func(FrameInfo, error) bool) {
	capture, err := gocv.OpenVideoCapture(v.path)
	if err != nil {
		yield(FrameInfo{}, err)
		return
	}
	defer capture.Close()

	modulo := v.moduloToAchieveTargetFps()

	currentFrame := gocv.NewMat()
	defer currentFrame.Close()
	for {
		if !capture.Read(&currentFrame) {
			break
		}

		frameNumber := int(capture.Get(gocv.VideoCapturePosFrames))
		if frameNumber%modulo != 0 {
			continue
		}

		info := FrameInfo{
			FrameNum:      frameNumber,
			Frame:         &currentFrame,
			TimeFromStart: v.currentTime(frameNumber),
		}
		if !yield(info, nil) {
			break
		}
	}
}

func (v VideoIterator) IterateWithPrevious(yield func(withPrevious FrameInfoWithPrevious, err error) bool) {
	previousMat := gocv.NewMat()
	previous := FrameInfo{
		Frame: &previousMat,
	}
	defer previous.Frame.Close()

	isFirst := true

	v.Iterate(func(fi FrameInfo, err error) bool {
		if err != nil {
			return yield(FrameInfoWithPrevious{}, err)
		}

		if isFirst {
			isFirst = false
			previous.FrameNum = fi.FrameNum
			previous.TimeFromStart = fi.TimeFromStart
			fi.Frame.CopyTo(previous.Frame)
			return true
		}

		if !yield(FrameInfoWithPrevious{Current: fi, Previous: previous}, nil) {
			return false
		}

		previous.FrameNum = fi.FrameNum
		previous.TimeFromStart = fi.TimeFromStart
		fi.Frame.CopyTo(previous.Frame)

		return true
	})
}
