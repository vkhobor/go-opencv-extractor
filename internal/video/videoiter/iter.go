package videoiter

import (
	"iter"
	"math"
	"time"

	"gocv.io/x/gocv"
)

type FrameInfo struct {
	Frame         gocv.Mat
	FrameNum      int
	TimeFromStart time.Duration
}

func (fi FrameInfo) Clone() FrameInfo {
	return FrameInfo{
		Frame:         fi.Frame.Clone(),
		FrameNum:      fi.FrameNum,
		TimeFromStart: fi.TimeFromStart,
	}
}

func moduloToAchieveTargetFps(originalFPS, targetFPS float64) int {
	return int(math.Ceil(float64(originalFPS) / float64(targetFPS)))
}

func AllSampledFrames(v Video, fpsWant int, progress func(Progress)) iter.Seq2[FrameInfo, error] {
	return func(yield func(FrameInfo, error) bool) {
		capture, err := gocv.OpenVideoCapture(v.path)
		if err != nil {
			yield(FrameInfo{}, err)
			return
		}
		defer capture.Close()

		capture.Set(gocv.VideoCapturePosFrames, float64(v.startFrame))

		currentFrame := gocv.NewMat()
		defer currentFrame.Close()

		samplingFactor := moduloToAchieveTargetFps(v.originalFPS, float64(fpsWant))

		for {
			capture.Grab(samplingFactor - 1)

			if !capture.Read(&currentFrame) || v.endFrame <= int(capture.Get(gocv.VideoCapturePosFrames)) {
				break
			}

			frameNumber := int(capture.Get(gocv.VideoCapturePosFrames))
			progress(v.CurrentProgress(frameNumber))

			info := FrameInfo{
				FrameNum:      frameNumber,
				Frame:         currentFrame,
				TimeFromStart: v.CurrentTime(frameNumber),
			}

			if !yield(info, nil) {
				break
			}
		}
	}
}
