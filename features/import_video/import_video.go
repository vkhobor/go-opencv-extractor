package import_video

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/features"
	"github.com/vkhobor/go-opencv/mlog"
	"github.com/vkhobor/go-opencv/video/videoiter"
	"gocv.io/x/gocv"

	"github.com/vkhobor/go-opencv/filters"
	"github.com/vkhobor/go-opencv/image/surf"
)

type ImportVideoFeature struct {
	SqlDB   features.TXer
	Querier features.QuerierWithTx
	Config  config.DirectoryConfig
}

func (d *ImportVideoFeature) ImportVideo(ctx context.Context, videoID string, jobID string, filterID string) error {
	tx, err := d.SqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	queries := d.Querier.WithTx(tx)

	// TODO make this more efficient, no need to query for every video
	downloadedVideos, err := queries.GetVideosDownloadedButNotImported(ctx)
	if err != nil {
		slog.Error("GetDownloadedVideos: Error while getting downloaded videos", "error", err)
		return err
	}
	if len(downloadedVideos) == 0 {
		return errors.New("no downloaded videos")
	}

	var videoSavePath string
	for _, video := range downloadedVideos {
		if video.YtVideoID == videoID {
			videoSavePath = video.Path
			break
		}
	}

	refs, err := d.getRefImages(ctx, tx, jobID)
	if err != nil {
		return err
	}
	if len(refs.Paths) == 0 {
		return errors.New("no ref images found")
	}

	id, err := d.startImportAttempt(ctx, tx, videoID, filterID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	videoError := d.handleSingle(ctx, id, refs, videoID, videoSavePath)

	tx, err = d.SqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if videoError != nil {
		innerErr := d.updateError(ctx, tx, id, err)
		if innerErr != nil {
			return err
		}
		return err
	} else {
		err = d.finishImport(ctx, tx, videoID, id)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (d *ImportVideoFeature) handleSingle(
	ctx context.Context,
	importAttemptId string,
	refs FilterWithPaths,
	videoID string,
	videoSavePath string) error {

	progress := make(chan videoiter.Progress, 1)
	defer close(progress)
	go d.importProgressHandler(ctx, progress, videoID, importAttemptId)
	progressHandler := func(p videoiter.Progress) {
		progress <- p
	}

	options := []surf.SURFImageMatcherOption{
		surf.WithMinMatches(int(refs.MinSURFMatches)),
		surf.WithMinThreshold(refs.MinThresholdForSURFMatches),
		surf.WithRatioThreshold(refs.RatioTestThreshold),
	}
	mlog.Log().Info("Initializing SURF image matcher with options", "minThreshold", refs.MinThresholdForSURFMatches, "minMatches", refs.MinSURFMatches, "ratioThreshold", refs.RatioTestThreshold)

	matcher, err := surf.NewSURFImageMatcher(refs.Paths[:], options...)
	if err != nil {
		return err
	}
	defer matcher.Close()

	mlog.Log().Info("Using mseskip", "mseskip", refs.MSESkip)
	filter := filters.NewSURFVideoFilter(matcher, refs.MSESkip)

	video, err := videoiter.NewVideo(videoSavePath)
	if err != nil {
		return err
	}

	fpsWant := filter.SamplingWantFPS()

	frames := videoiter.AllSampledFrames(video, fpsWant, progressHandler)
	wantFrames := filter.FrameFilter(frames)

	err = d.collectFramesToDisk(ctx, wantFrames, d.Config.GetImagesDir(), videoID, importAttemptId)
	if err != nil {
		return err
	}

	slog.Info("Processed images", "outputDir", d.Config.GetImagesDir(), "fpsWant", fpsWant)

	mlog.Log().Info("Imported video", "id", videoID)
	return nil
}

func (d *ImportVideoFeature) importProgressHandler(
	ctx context.Context,
	progressChan <-chan videoiter.Progress,
	videoID string,
	importAttemptId string) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	previous := videoiter.Progress{}
	for item := range progressChan {
		select {
		case <-ticker.C:
			err := d.updateProgress(ctx, importAttemptId, int(item.Percent()))
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

func (d *ImportVideoFeature) collectFramesToDisk(ctx context.Context, frames iter.Seq2[videoiter.FrameInfo, error], outputDir string, videoId string, importAttemptId string) error {
	for value, err := range frames {
		if err != nil {
			return err
		}

		filePath, ok := saveFrameWithUUIDName(outputDir, &value.Frame)
		if !ok {
			return errors.New("failed to save frame")
		}

		d.addFrameToVideo(ctx, Frame{
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
