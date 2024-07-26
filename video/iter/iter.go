package iter

import (
	"errors"
	"fmt"
	"io"
	"iter"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/vkhobor/go-opencv/video/metadata"
	"gocv.io/x/gocv"
)

type Progress struct {
	Done      int
	Max       int
	timestamp time.Time
}

func (p Progress) Percent() float64 {
	return float64(p.Done) / float64(p.Max)
}

func (p Progress) FPS(other Progress) float64 {
	return math.Abs(float64(p.Done-other.Done)) / p.timestamp.Sub(other.timestamp).Seconds()
}

func (p Progress) MergeWith(other ...Progress) Progress {
	for _, o := range other {
		p.Done += o.Done
		p.Max += o.Max
		if o.timestamp.After(p.timestamp) {
			p.timestamp = o.timestamp
		}

	}
	return p
}

type Video struct {
	path        string
	bufferSize  int
	sampleRate  int
	startFrame  int
	endFrame    int
	originalFPS float64
	progress    func(Progress)
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

func moduloToAchieveTargetFps(originalFPS, targetFPS float64) int {
	return int(math.Ceil(float64(originalFPS) / float64(targetFPS)))
}

func NewVideo(path string, fpsWant int, bufferSize int, progress func(Progress)) (Video, error) {
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

	samplingFactor := moduloToAchieveTargetFps(meta, float64(fpsWant))

	iter := Video{
		startFrame:  0,
		bufferSize:  bufferSize,
		endFrame:    maxFrame,
		path:        path,
		originalFPS: meta,
		sampleRate:  samplingFactor,
		progress:    progress,
	}

	return iter, nil
}

func copyFile(src, dst string) {
	fin, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer fin.Close()

	fout, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	_, err = io.Copy(fout, fin)

	if err != nil {
		panic(err)
	}
}
func MergeSliceProgress(progresses []Progress) Progress {
	p := Progress{}
	for _, pr := range progresses {
		p = p.MergeWith(pr)
	}
	return p
}

func (v Video) SplitToChunks(chunks int) []Video {
	videoChunks := make([]Video, chunks)
	progresses := make([]Progress, chunks)
	for i := range chunks {
		start := v.startFrame + i*(v.MaxFrame()/chunks)
		end := start + v.MaxFrame()/chunks
		if i == chunks-1 {
			end = v.endFrame
		}

		dir := filepath.Dir(v.path)
		filename := filepath.Base(v.path)
		newPath := filepath.Join(dir, fmt.Sprintf("%s_%d.mp4", filename, i))
		copyFile(v.path, newPath)
		progresses[i] = Progress{
			Done:      0,
			Max:       v.MaxFrame(),
			timestamp: time.Now(),
		}

		videoChunks[i] = Video{
			startFrame: start,
			endFrame:   end,
			path:       newPath,
			progress: func(p Progress) {
				progresses[i] = p
				merged := MergeSliceProgress(progresses)
				v.progress(merged)
			},
			originalFPS: v.originalFPS,
			sampleRate:  v.sampleRate,
		}
	}
	return videoChunks
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

func (v Video) currentTime(frame int) time.Duration {
	return time.Second * time.Duration(math.Ceil(float64(frame)/v.originalFPS))
}

func (v Video) CurrentProgress(frame int) Progress {
	return Progress{
		Done:      frame - v.startFrame,
		Max:       v.MaxFrame(),
		timestamp: time.Now(),
	}
}

func (v Video) AllSampledFrames() iter.Seq2[FrameInfo, error] {
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
		for {
			capture.Grab(v.sampleRate - 1)

			if !capture.Read(&currentFrame) || v.endFrame <= int(capture.Get(gocv.VideoCapturePosFrames)) {
				break
			}

			frameNumber := int(capture.Get(gocv.VideoCapturePosFrames))
			v.progress(v.CurrentProgress(frameNumber))
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
}

func (v Video) BufferedFrames() iter.Seq2[FrameInfo, error] {
	return func(yield func(FrameInfo, error) bool) {
		buffer := make([]FrameInfo, v.bufferSize)
		defer func() {
			for _, frame := range buffer {
				frame.Frame.Close()
			}
		}()

		for fi, err := range v.AllSampledFrames() {
			frameCopy := fi.Frame.Clone()
			if len(buffer) < v.bufferSize {
				buffer = append(buffer, FrameInfo{
					Frame:         &frameCopy,
					FrameNum:      fi.FrameNum,
					TimeFromStart: fi.TimeFromStart,
				})
				continue
			}

			for _, frame := range buffer {
				if !yield(frame, err) {
					break
				}
			}

			for _, frame := range buffer {
				frame.Frame.Close()
			}
			buffer = make([]FrameInfo, v.bufferSize)
			buffer = append(buffer, FrameInfo{
				Frame:         &frameCopy,
				FrameNum:      fi.FrameNum,
				TimeFromStart: fi.TimeFromStart,
			})
		}
	}
}

func (v Video) AllFramesWithPrevious() iter.Seq2[FrameInfoWithPrevious, error] {
	return func(yield func(FrameInfoWithPrevious, error) bool) {
		previousMat := gocv.NewMat()
		previous := FrameInfo{
			Frame: &previousMat,
		}
		defer previous.Frame.Close()

		isFirst := true

		for fi, err := range v.AllSampledFrames() {
			if err != nil {
				cont := yield(FrameInfoWithPrevious{}, err)
				if !cont {
					break
				}
				continue
			}

			if isFirst {
				isFirst = false
				previous.FrameNum = fi.FrameNum
				previous.TimeFromStart = fi.TimeFromStart
				fi.Frame.CopyTo(previous.Frame)
				continue
			}

			if !yield(FrameInfoWithPrevious{Current: fi, Previous: previous}, nil) {
				break
			}

			previous.FrameNum = fi.FrameNum
			previous.TimeFromStart = fi.TimeFromStart
			fi.Frame.CopyTo(previous.Frame)
		}
	}
}
