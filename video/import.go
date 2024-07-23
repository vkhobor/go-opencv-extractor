package video

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/vkhobor/go-opencv/image"
	"github.com/vkhobor/go-opencv/iter"
	"github.com/vkhobor/go-opencv/mlog"
	videoiter "github.com/vkhobor/go-opencv/video/iter"

	"github.com/google/uuid"
	"gocv.io/x/gocv"
)

func HandleVideoFromPath(path string, outputDir string, fpsWant int, refImagePaths []string, progress func(float64)) ([]string, error) {
	iterator, err := videoiter.NewVideoIterator(path, fpsWant)
	if err != nil {
		return nil, err
	}
	mlog.Log().Debug("VideoIterator created", "path", path, "fps", fpsWant)

	checker, err := image.NewChecker(refImagePaths)
	if err != nil {
		return nil, err
	}
	defer checker.Close()

	mlog.Log().Debug("Checker created", "refImagePaths", refImagePaths)

	surfMatch := iter.Filter2(iterator.IterateWithPrevious, func(info videoiter.FrameInfoWithPrevious, err error) bool {
		progress(float64(info.Current.FrameNum) / float64(iterator.MaxFrame()) * 100)
		if err != nil {
			return true
		}

		return checker.IsImageMatch(*info.Current.Frame)
	})

	surfMatchEnoughDifference := iter.FilterCanError(surfMatch, func(info videoiter.FrameInfoWithPrevious, err error) (bool, error) {
		if err != nil {
			return false, err
		}
		if distanceIsLessThanDuration(info.Current, info.Previous, time.Minute*2) {
			diff, err := image.CompareImages(info.Previous.Frame, info.Current.Frame)
			if err != nil {
				return true, err
			}

			// if diff is too small, skip
			if diff < 0.2 {
				return false, nil
			}
		}
		return true, nil
	})

	filePaths := make([]string, 0)
	var iterationError error
	surfMatchEnoughDifference(func(value videoiter.FrameInfoWithPrevious, err error) bool {
		if err != nil {
			iterationError = err
			return false
		}

		filePath, ok := saveFrameWithUUIDName(outputDir, value.Current)
		if !ok {
			iterationError = fmt.Errorf("failed to write image to file: %v", filePath)
			return false
		}
		filePaths = append(filePaths, filePath)

		return true
	})

	if iterationError != nil {
		return nil, iterationError
	}

	slog.Info("Processed images", "fileNames", filePaths, "outputDir", outputDir, "fpsWant", fpsWant, "count", len(filePaths))
	return filePaths, nil
}

func saveFrameWithUUIDName(outputDir string, value videoiter.FrameInfo) (string, bool) {
	fileName := fmt.Sprintf("%v.jpg", uuid.New().String())
	filePath := filepath.Join(outputDir, fileName)
	ok := gocv.IMWrite(filePath, *value.Frame)
	return filePath, ok
}

func distanceIsLessThanDuration(frame1 videoiter.FrameInfo, frame2 videoiter.FrameInfo, duration time.Duration) bool {
	return frame2.TimeFromStart-frame1.TimeFromStart < duration
}
