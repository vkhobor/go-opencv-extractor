package video

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/vkhobor/go-opencv/image"
	"github.com/vkhobor/go-opencv/iter"
	"github.com/vkhobor/go-opencv/mlog"
	videoiter "github.com/vkhobor/go-opencv/video/iter"

	"github.com/google/uuid"
	"gocv.io/x/gocv"
)

type result struct {
	Error  error
	Result []string
}

func HandleVideoFromPath(
	path string,
	outputDir string,
	fpsWant int,
	refImagePaths []string,
	progress func(p videoiter.Progress)) ([]string, error) {

	video, err := videoiter.NewVideo(path, fpsWant, 500, progress)
	if err != nil {
		return nil, err
	}

	chunkLength := runtime.GOMAXPROCS(0)
	chunkLength = 1
	chunks := video.SplitToChunks(chunkLength)

	mlog.Log().Debug("Video created",
		"path", path,
		"fps", fpsWant,
		"chunks", chunks,
		"maxFrame", video.MaxFrame(),
		"chunkLength", chunkLength)

	wg := sync.WaitGroup{}
	wg.Add(chunkLength)

	var results = make([]result, chunkLength)
	for i := range results {
		results[i] = result{Error: nil, Result: []string{}}
	}

	for index, chunk := range chunks {
		mlog.Log().Debug("Processing chunk", "path", chunk.Path())
		go func() {
			singleResult, err := handleVideoIter(refImagePaths, chunk, fpsWant, outputDir)
			results[index] = result{Error: err, Result: singleResult}
			wg.Done()
		}()
	}
	wg.Wait()

	return collectResults(results[:])
}

func collectResults(results []result) ([]string, error) {
	errs := make([]error, 0)
	filePaths := make([]string, 0)
	for _, result := range results {
		if result.Error != nil {
			errs = append(errs, result.Error)
		} else {
			filePaths = append(filePaths, result.Result...)
		}
	}
	return filePaths, errors.Join(errs...)
}

func handleVideoIter(refImagePaths []string, video videoiter.Video, fpsWant int, outputDir string) ([]string, error) {
	checker, err := image.NewChecker(refImagePaths)
	if err != nil {
		return nil, err
	}
	defer checker.Close()
	mlog.Log().Debug("Checker created", "refImagePaths", refImagePaths)

	frames := video.AllFramesWithPrevious()

	frames = iter.FilterWithError2(
		frames,
		func(info videoiter.FrameInfoWithPrevious, err error) (bool, error) {
			if err == nil {
				if distanceIsLessThanDuration(info.Previous, info.Current, time.Minute*2) {
					diff, err := image.CompareImages(info.Previous.Frame, info.Current.Frame)
					if err != nil {
						return true, err
					}

					if diff < 0.2 {
						return false, nil
					}
				}
				return true, nil
			}

			return true, err
		})

	frames = iter.Filter2(
		frames,
		func(info videoiter.FrameInfoWithPrevious, err error) bool {
			if err == nil {
				return checker.IsImageMatch(*info.Current.Frame)
			}
			return true
		})

	filePaths := make([]string, 0)
	var iterationError error

	for value, err := range frames {
		if err != nil {
			iterationError = err
			break
		}

		filePath, ok := saveFrameWithUUIDName(outputDir, value.Current.Frame)
		if !ok {
			iterationError = errors.New("failed to save frame")
			break
		}
		filePaths = append(filePaths, filePath)
	}

	if iterationError != nil {
		return nil, iterationError
	}

	slog.Info("Processed images", "fileNames", filePaths, "outputDir", outputDir, "fpsWant", fpsWant, "count", len(filePaths))
	return filePaths, nil
}

func saveFrameWithUUIDName(outputDir string, value *gocv.Mat) (string, bool) {
	fileName := fmt.Sprintf("%v.jpg", uuid.New().String())
	filePath := filepath.Join(outputDir, fileName)
	mlog.Log().Debug("Saving frame", "filePath", filePath)
	ok := gocv.IMWrite(filePath, *value)
	return filePath, ok
}

func distanceIsLessThanDuration(one, two videoiter.FrameInfo, duration time.Duration) bool {
	dist := one.TimeFromStart - two.TimeFromStart
	return dist.Abs() < duration
}
