package background

import (
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
	Output   chan<- queries.ImportedVideo
	Config   config.DirectoryConfig
}

func (d *Importer) Start() {
	for video := range d.Input {
		imported, err := d.importVideo(video)
		if err != nil {
			slog.Error("Error while importing video", "error", err, "video", video)
			continue
		}

		d.Output <- imported
	}
}

func (d *Importer) importVideo(video queries.DownlodedVideo) (queries.ImportedVideo, error) {
	refs, err := d.Queries.GetRefImages(video)
	if err != nil || len(refs) == 0 {
		slog.Debug("No videos to import or no reference images", "video", video, "refs", refs, "error", err)
		return queries.ImportedVideo{}, nil
	}

	id, err := d.Queries.StartImportAttempt(video)
	if err != nil {
		slog.Error("Error while starting import attempt", "error", err, "video", video)
		return queries.ImportedVideo{}, err
	}

	slog.Info("Running import job", "videos_waiting_for_import", d)
	videoImported := d.handleSingle(refs, video)

	// TODO should only allow saving if does not make the database inconsistent, like multiple sucessful imports for the same video
	d.Queries.SaveFrames(videoImported, id)
	d.Queries.UpdateProgress(id, 100)
	return videoImported, nil
}

func (d *Importer) handleSingle(refs []string, video queries.DownlodedVideo) queries.ImportedVideo {
	time.Sleep(d.Throttle)

	// Make progress handling async
	progress := make(chan float64, 1)
	defer close(progress)
	go d.importProgressHandler(progress, video)
	progressHandler := func(progressFromImport float64) {
		progress <- progressFromImport
	}

	val, err := videoLib.HandleVideoFromPath(video.SavePath, d.Config.GetImagesDir(), 1, "", refs, progressHandler)
	if err != nil {
		slog.Error("Error while importing video", "error", err, "video", video)
		return queries.ImportedVideo{
			Error: err,
		}
	}
	slog.Info("Imported video", "video", val.FilePaths, "id", video.ID)
	frames := make([]queries.Frame, 0)
	for _, v := range val.FilePaths {
		frames = append(frames, queries.Frame{FrameNumber: 0, Path: v})
	}
	return queries.ImportedVideo{
		DownlodedVideo:  video,
		ExtractedFrames: frames,
	}
}

func (d *Importer) importProgressHandler(progressChan <-chan float64, video queries.DownlodedVideo) {
	// sample progresschan every one 30 seconds
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for item := range progressChan {
		select {
		case <-ticker.C:
			d.Queries.UpdateProgress(video.ID, int(item*100))
			slog.Info("importProgressHandler", "id", video.ID, "progress", item)
		default:
		}
	}
}
