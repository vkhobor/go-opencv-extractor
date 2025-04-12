package features

import (
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/filters"
	"github.com/vkhobor/go-opencv/image/surf"
	"github.com/vkhobor/go-opencv/mlog"
	"github.com/vkhobor/go-opencv/queries"
	"github.com/vkhobor/go-opencv/video/videoiter"
	"gocv.io/x/gocv"
)

type ImportVideoFeature struct {
	Queries *queries.Queries
	Config  config.DirectoryConfig
}

func (d *ImportVideoFeature) ImportVideo(videoID string, jobID string, filterID string) error {
	// TODO check if video is already imported, optionally abort while progressing
	// TODO make this more efficient, no need to query for every video
	downloadedVideos := d.Queries.GetDownloadedVideos(false)
	if len(downloadedVideos) == 0 {
		return errors.New("no downloaded videos")
	}

	var videoSavePath string
	for _, video := range downloadedVideos {
		if video.ID == videoID {
			videoSavePath = video.SavePath
			break
		}
	}

	refs, err := d.Queries.GetRefImages(jobID)
	if err != nil {
		return err
	}
	if len(refs) == 0 {
		return errors.New("no ref images found")
	}

	id, err := d.Queries.StartImportAttempt(videoID, filterID)
	if err != nil {
		return err
	}

	err = d.handleSingle(id, refs, videoID, videoSavePath)
	if err != nil {
		innerErr := d.Queries.UpdateError(id, err)
		if innerErr != nil {
			return err
		}
		return err
	}

	err = d.Queries.FinishImport(videoID, id)
	if err != nil {
		if errors.Is(err, queries.ErrHasImported) {
			innerErr := d.Queries.UpdateError(id, err)
			if innerErr != nil {
				return err
			}
		}

		return err
	}

	return nil
}

func (d *ImportVideoFeature) handleSingle(
	importAttemptId string,
	refs []string,
	videoID string,
	videoSavePath string) error {

	// Make progress handling async
	progress := make(chan videoiter.Progress, 1)
	defer close(progress)
	go d.importProgressHandler(progress, videoID, importAttemptId)
	progressHandler := func(p videoiter.Progress) {
		progress <- p
	}

	// TODO map different filters if implemented in database
	matcher, err := surf.NewSURFImageMatcher(refs)
	if err != nil {
		return err
	}
	defer matcher.Close()

	filter := filters.NewSURFVideoFilter(matcher)

	video, err := videoiter.NewVideo(videoSavePath)
	if err != nil {
		return err
	}

	fpsWant := filter.SamplingWantFPS()

	frames := videoiter.AllSampledFrames(video, fpsWant, progressHandler)
	wantFrames := filter.FrameFilter(frames)

	err = d.collectFramesToDisk(wantFrames, d.Config.GetImagesDir(), videoID, importAttemptId)
	if err != nil {
		return err
	}

	slog.Info("Processed images", "outputDir", d.Config.GetImagesDir(), "fpsWant", fpsWant)

	mlog.Log().Info("Imported video", "id", videoID)
	return nil
}

func (d *ImportVideoFeature) importProgressHandler(
	progressChan <-chan videoiter.Progress,
	videoID string,
	importAttemptId string) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	previous := videoiter.Progress{}
	for item := range progressChan {
		select {
		case <-ticker.C:
			err := d.Queries.UpdateProgress(importAttemptId, int(item.Percent()))
			if err != nil {
				mlog.Log().Error("Failed to update progress", "error", err)
			}
			mlog.Log().Info("importProgressHandler",
				"id", videoID,
				"progress", item,
				"speed fps", item.FPS(previous),
				"percent", item.Percent())
			previous = item
		default:
		}
	}
}

type Filter interface {
	FrameFilter(frameSeq iter.Seq2[videoiter.FrameInfo, error]) iter.Seq2[videoiter.FrameInfo, error]
	SamplingWantFPS() int
}

func (d *ImportVideoFeature) collectFramesToDisk(frames iter.Seq2[videoiter.FrameInfo, error], outputDir string, videoId string, importAttemptId string) error {
	for value, err := range frames {
		if err != nil {
			return err
		}

		filePath, ok := saveFrameWithUUIDName(outputDir, &value.Frame)
		if !ok {
			return errors.New("failed to save frame")
		}

		d.Queries.AddFrameToVideo(videoId, queries.Frame{
			FrameNumber: value.FrameNum,
			Path:        filePath,
		}, importAttemptId)
		mlog.Log().Info("Saving frame", "filePath", filePath)
	}

	return nil
}

func saveFrameWithUUIDName(outputDir string, value *gocv.Mat) (string, bool) {
	fileName := fmt.Sprintf("%v.jpg", uuid.New().String())
	filePath := filepath.Join(outputDir, fileName)
	mlog.Log().Debug("Saving frame", "filePath", filePath)
	ok := gocv.IMWrite(filePath, *value)
	return filePath, ok
}
