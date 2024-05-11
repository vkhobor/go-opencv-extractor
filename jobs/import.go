package jobs

import (
	"log/slog"
	"time"

	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/domain"
	videoLib "github.com/vkhobor/go-opencv/video"
)

type Importer struct {
	queries      *domain.JobQueries
	importSingle func(refs []string, video domain.DownlodedVideo) domain.ImportedVideo
}

// NewImporter creates a new Importer
func NewImporter(queries *domain.JobQueries, throttle time.Duration, config config.DirectoryConfig) *Importer {

	handleSingle := func(refs []string, video domain.DownlodedVideo) domain.ImportedVideo {
		time.Sleep(throttle)

		progress := make(chan float64)
		defer close(progress)
		go importProgressHandler(progress, video)

		val, err := videoLib.HandleVideoFromPath(video.SavePath, config.GetImagesDir(), 1, "", refs, progress)
		if err != nil {
			slog.Error("Error while importing video", "error", err, "video", video)
			return domain.ImportedVideo{
				Error: err,
			}
		}
		slog.Info("Imported video", "video", val.FilePaths, "id", video.ID)
		frames := make([]domain.Frame, 0)
		for _, v := range val.FilePaths {
			frames = append(frames, domain.Frame{FrameNumber: 0, Path: v})
		}
		return domain.ImportedVideo{
			DownlodedVideo:  video,
			ExtractedFrames: frames,
		}
	}

	return &Importer{
		queries:      queries,
		importSingle: handleSingle,
	}
}

// Imports all importable from db
func (jc *Importer) ImportAllImportableFromDb() {

	d := jc.queries.GetDownloadedVideos()
	refs, err := jc.queries.GetRefImages()
	if err != nil || len(refs) == 0 || len(d) == 0 {
		slog.Debug("No videos to import or no reference images")
		return
	}

	slog.Info("Running import job", "videos_waiting_for_import", d)
	for _, v := range d {
		video := jc.importSingle(refs, v)
		jc.queries.SaveImported(video)
	}
}

func importProgressHandler(progressChan <-chan float64, video domain.DownlodedVideo) {
	// sample progresschan every one 30 seconds
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for item := range progressChan {
		select {
		case <-ticker.C:
			slog.Info("Progress", "id", video.ID, "progress", item)
		default:
		}
	}
}
