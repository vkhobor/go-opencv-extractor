package importing

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/vkhobor/go-opencv/image"

	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"
	"gocv.io/x/gocv"
)

type ImportContext struct {
	url    string
	output string
	fps    int
}

func HandleVideoFromPath(path string, outputDir string, fpsWant int, videoTitle string, refImagePaths []string) (*DbEntry, error) {
	fps, err := extractMetadata(path)
	if err != nil {
		return nil, err
	}

	progressChan := make(chan int)
	defer close(progressChan)

	iter, err := image.NewExtractIterator(image.Config{
		VideoPath:        path,
		PathsToRefImages: refImagePaths,
		OriginalFPS:      fps,
		WantFPS:          fpsWant,
	}, progressChan)

	if err != nil {
		return nil, err
	}
	defer iter.Close()

	go func() {
		bar := progressbar.Default(int64(iter.Length()))
		defer bar.Finish()

		for progress := range progressChan {
			bar.Set(progress)
		}
	}()

	fileNames, err := processImages(iter, progressChan, outputDir, fpsWant)
	if err != nil {
		return nil, err
	}

	return &DbEntry{Status: StatusImported, Title: videoTitle, FileNames: fileNames}, nil
}

func processImages(iter *image.ExtractIterator, progressChan chan<- int, outputDir string, fpsWant int) ([]string, error) {
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

	fmt.Printf("Found %v\n", len(fileNames))
	return fileNames, nil
}
