package background

import (
	"errors"
	"log/slog"
	"time"

	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/queries"
	videoLib "github.com/vkhobor/go-opencv/video"
)

type Importer struct {
	Queries  *queries.Queries
	Throttle time.Duration
	Input    <-chan queries.DownlodedVideo
	Config   config.DirectoryConfig
}

func (d *Importer) Start() {
	for video := range d.Input {
		slog.Debug("Importer starting importing", "video", video, "method", "Start")
		_, err := d.importVideo(video)
		if err != nil {
			slog.Error("Error while importing video", "error", err, "video", video)
			continue
		}

		// TODO handle output if needed
		// d.Output <- imported
	}
}

// TODO optionally move the single processing to another package e.g import/service
func (d *Importer) importVideo(video queries.DownlodedVideo) (queries.ImportedVideo, error) {
	refs, err := d.Queries.GetRefImages(video)
	if err != nil {
		return queries.ImportedVideo{}, err
	}
	if len(refs) == 0 {
		return queries.ImportedVideo{}, errors.New("no ref images found")
	}

	id, err := d.Queries.StartImportAttempt(video)
	if err != nil {
		return queries.ImportedVideo{}, err
	}

	videoImported, err := d.handleSingle(id, refs, video)
	if err != nil {
		innerErr := d.Queries.UpdateError(id, err)
		if innerErr != nil {
			return queries.ImportedVideo{}, err
		}
		return queries.ImportedVideo{}, err
	}

	err = d.Queries.FinishImport(videoImported, id)
	if err != nil {
		if errors.Is(err, queries.ErrHasImported) {
			innerErr := d.Queries.UpdateError(id, err)
			if innerErr != nil {
				return queries.ImportedVideo{}, err
			}
		}

		return queries.ImportedVideo{}, err
	}

	return videoImported, nil
}

func (d *Importer) handleSingle(importAttemptId string, refs []string, video queries.DownlodedVideo) (queries.ImportedVideo, error) {
	time.Sleep(d.Throttle)

	// Make progress handling async
	progress := make(chan float64, 1)
	defer close(progress)
	go d.importProgressHandler(progress, video, importAttemptId)
	progressHandler := func(progressFromImport float64) {
		progress <- progressFromImport
	}

	val, err := videoLib.HandleVideoFromPath(video.SavePath, d.Config.GetImagesDir(), 1, "", refs, progressHandler)
	if err != nil {
		return queries.ImportedVideo{}, err
	}

	slog.Info("Imported video", "video", val.FilePaths, "id", video.ID)
	frames := make([]queries.Frame, 0)
	for _, v := range val.FilePaths {
		frames = append(frames, queries.Frame{FrameNumber: 0, Path: v})
	}
	return queries.ImportedVideo{
		DownlodedVideo:  video,
		ExtractedFrames: frames,
	}, nil
}

func (d *Importer) importProgressHandler(progressChan <-chan float64, video queries.DownlodedVideo, importAttemptId string) {
	// sample progresschan every one 30 seconds
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for item := range progressChan {
		select {
		case <-ticker.C:
			_ = d.Queries.UpdateProgress(importAttemptId, int(item*100))
			slog.Info("importProgressHandler", "id", video.ID, "progress", item)
		default:
		}
	}
}
