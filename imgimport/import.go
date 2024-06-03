package imgimport

import (
	"log/slog"
	"time"

	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/download"
	videoLib "github.com/vkhobor/go-opencv/video"
)

type Frame struct {
	FrameNumber int
	Path        string
}

type ImportedVideo struct {
	download.DownlodedVideo
	ExtractedFrames []Frame
	Error           error
}

type Importer struct {
	Queries  *Queries
	Throttle time.Duration
	Input    <-chan download.DownlodedVideo
	Output   chan<- ImportedVideo
	Config   config.DirectoryConfig
}

func (d *Importer) Start() {
	for video := range d.Input {
		imported, err := d.importVideo(video)
		if err != nil {
			slog.Error("Error while importing video", "error", err, "video", video)
		}
		d.Output <- imported
	}
}

func (d *Importer) importVideo(video download.DownlodedVideo) (ImportedVideo, error) {
	refs, err := d.Queries.GetRefImages(video)
	if err != nil || len(refs) == 0 {
		slog.Debug("No videos to import or no reference images")
		return ImportedVideo{}, nil
	}

	id, err := d.Queries.StartImportAttempt(video)
	if err != nil {
		slog.Error("Error while starting import attempt", "error", err, "video", video)
		return ImportedVideo{}, err
	}

	slog.Info("Running import job", "videos_waiting_for_import", d)
	videoImported := d.handleSingle(refs, video)

	d.Queries.UpdateProgress(id, 100)
	return videoImported, nil
}

func (d *Importer) handleSingle(refs []string, video download.DownlodedVideo) ImportedVideo {
	time.Sleep(d.Throttle)

	progress := make(chan float64)
	defer close(progress)
	go d.importProgressHandler(progress, video)

	val, err := videoLib.HandleVideoFromPath(video.SavePath, d.Config.GetImagesDir(), 1, "", refs, progress)
	if err != nil {
		slog.Error("Error while importing video", "error", err, "video", video)
		return ImportedVideo{
			Error: err,
		}
	}
	slog.Info("Imported video", "video", val.FilePaths, "id", video.ID)
	frames := make([]Frame, 0)
	for _, v := range val.FilePaths {
		frames = append(frames, Frame{FrameNumber: 0, Path: v})
	}
	return ImportedVideo{
		DownlodedVideo:  video,
		ExtractedFrames: frames,
	}
}

func (d *Importer) importProgressHandler(progressChan <-chan float64, video download.DownlodedVideo) {
	// sample progresschan every one 30 seconds
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for item := range progressChan {
		select {
		case <-ticker.C:
			d.Queries.UpdateProgress(video.ID, int(item*100))
			slog.Info("Progress", "id", video.ID, "progress", item)
		default:
		}
	}
}
