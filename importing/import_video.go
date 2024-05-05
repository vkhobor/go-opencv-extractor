package importing

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/vkhobor/go-opencv/image"

	"github.com/google/uuid"
	"gocv.io/x/gocv"
)

type ImportContext struct {
	url    string
	output string
	fps    int
}

func HandleVideoFromPath(path string, outputDir string, fpsWant int, videoTitle string, refImagePaths []string, progressChan chan<- float64) (*DbEntry, error) {
	fps, err := extractMetadata(path)
	if err != nil {
		return nil, err
	}

	frameChan := make(chan struct{})
	defer close(frameChan)

	iter, err := image.NewExtractIterator(image.Config{
		VideoPath:        path,
		PathsToRefImages: refImagePaths,
		OriginalFPS:      fps,
		WantFPS:          fpsWant,
	}, frameChan)

	if err != nil {
		return nil, err
	}
	defer iter.Close()

	go func() {
		for range frameChan {
			progressChan <- float64(iter.CurrentFrame()) / float64(iter.Length())
		}
	}()

	fileNames, err := processImages(iter, outputDir, fpsWant)
	if err != nil {
		return nil, err
	}

	return &DbEntry{Status: StatusImported, Title: videoTitle, FileNames: fileNames}, nil
}

func processImages(iter *image.ExtractIterator, outputDir string, fpsWant int) ([]string, error) {
	fileNames := make([]string, 0)

	prev := gocv.NewMat()
	var prevFrame int
	for iter.Next() {
		err := func() error {
			value := iter.Value()

			currFrame := iter.CurrentFrame()

			if prevFrame != 0 && (currFrame-prevFrame)/fpsWant*1000 < int(time.Minute*2) {
				diff, err := image.CompareImages(prev, value)
				if err != nil {
					return err
				}

				// if diff is too small, skip
				if diff < 0.2 {
					return nil
				}
			}

			fileName := fmt.Sprintf("%v.jpg", uuid.New().String())
			filePath := filepath.Join(outputDir, fileName)
			gocv.IMWrite(filePath, value)
			fileNames = append(fileNames, fileName)

			value.CopyTo(&prev)
			prevFrame = currFrame
			return nil
		}()
		if err != nil {
			return nil, err
		}
	}

	slog.Info("Processed images", "fileNames", fileNames, "outputDir", outputDir, "fpsWant", fpsWant, "count", len(fileNames))
	return fileNames, nil
}
