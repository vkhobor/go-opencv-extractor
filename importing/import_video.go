package importing

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/vkhobor/go-opencv/db"
	"github.com/vkhobor/go-opencv/image"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"
	"gocv.io/x/gocv"
)

type ImportContext struct {
	url    string
	output string
	fps    int
}

func importVideoDbWrapper(id string, db *db.Db[DbEntry], operation func() (*DbEntry, error)) error {
	ctx := context.Background()

	// Update the database
	for {
		ok, entry, err := db.Read(ctx, id)
		if err != nil {
			return err
		}

		if ok && (entry.Status == StatusImported || entry.Status == StatusImporting) {
			return errors.New("video has already been imported or is being imported")
		}

		err = db.TryPut(id, DbEntry{Status: StatusImporting})
		if err == nil {
			break
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second * 5):
			// try again
		}
	}

	res, err := operation()
	defer func() {
		if err := recover(); err != nil {
			_ = db.Put(ctx, id, DbEntry{Status: StatusError, ErrorString: fmt.Sprintf("%v", err)})
			fmt.Println(err)
		}
	}()

	if err != nil {
		_ = db.Put(ctx, id, DbEntry{Status: StatusError, ErrorString: err.Error()})
		return err
	}

	errPut := db.Put(ctx, id, *res)
	if errPut != nil {
		return errPut
	}

	return nil
}

func ImportVideo(url string, outputDir string, fpsWant int, db *db.Db[DbEntry]) error {
	id, err := youtubeParser(url)
	if err != nil {
		return err
	}

	return importVideoDbWrapper(id, db, func() (*DbEntry, error) {
		return handleVideo(url, outputDir, fpsWant)
	})
}

func ImportVideoFromPath(id string, path string, outputDir string, fpsWant int, db *db.Db[DbEntry]) error {
	return importVideoDbWrapper(id, db, func() (*DbEntry, error) {
		return HandleVideoFromPath(path, outputDir, fpsWant, id)
	})
}

func HandleVideoFromPath(path string, outputDir string, fpsWant int, videoTitle string) (*DbEntry, error) {
	fps, err := extractMetadata(path)
	if err != nil {
		return nil, err
	}

	progressChan := make(chan int)
	defer close(progressChan)

	iter, err := image.NewExtractIterator(image.Config{
		VideoPath:   path,
		OriginalFPS: fps,
		WantFPS:     fpsWant,
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

func handleVideo(url string, outputDir string, fpsWant int) (*DbEntry, error) {
	videoPath, videoTitle, err := DownloadVideo(url)
	if err != nil {
		return nil, err
	}

	return HandleVideoFromPath(videoPath, outputDir, fpsWant, videoTitle)
}
