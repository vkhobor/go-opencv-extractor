package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/db_sql"
	"github.com/vkhobor/go-opencv/importing"
	"github.com/vkhobor/go-opencv/scraper"
)

type JobManager struct {
	Queries          *db_sql.Queries
	Wake             <-chan struct{}
	AutoWakePeriod   time.Duration
	ScrapeThrottle   time.Duration
	ImportThrottle   time.Duration
	DownloadThrottle time.Duration
}

func ImportProgressHandler(progressChan <-chan float64, video DownlodedVideo) {
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

func (jm *JobManager) Run() {

	scraper := func(args ScrapeArgs, ctx context.Context) <-chan ScrapedVideo {
		return MapChannel(scraper.ScrapeToChannel(args.SearchQuery, ctx, jm.ScrapeThrottle), func(id string) ScrapedVideo {
			slog.Info("Scraped and found video", "id", id)
			return ScrapedVideo{ID: id}
		})
	}

	importer := func(refs []string, vid ...DownlodedVideo) <-chan ImportedVideo {
		output := make(chan ImportedVideo)
		go func() {
			for _, video := range vid {
				time.Sleep(jm.ImportThrottle)

				progress := make(chan float64)
				defer close(progress)
				go ImportProgressHandler(progress, video)

				val, err := importing.HandleVideoFromPath(video.SavePath, config.WorkDirImages, 1, "", refs, progress)
				if err != nil {
					slog.Error("Error while importing video", "error", err, "video", video)
					output <- ImportedVideo{
						Error: err,
					}
					continue
				}
				slog.Info("Imported video", "video", val.FilePaths, "id", video.ID)
				frames := make([]Frame, 0)
				for _, v := range val.FilePaths {
					frames = append(frames, Frame{FrameNumber: 0, Path: v})
				}
				output <- ImportedVideo{
					DownlodedVideo:  video,
					ExtractedFrames: frames,
				}
			}
			close(output)
		}()
		return output
	}

	dowloader := func(vid ...ScrapedVideo) <-chan DownlodedVideo {
		output := make(chan DownlodedVideo)
		go func() {
			for _, video := range vid {
				time.Sleep(jm.DownloadThrottle)
				slog.Info("Download started", "video", video)
				path, title, err := importing.DownloadVideo(video.ID)
				if err != nil {
					output <- DownlodedVideo{
						ScrapedVideo: video,
						Error:        err,
					}
					slog.Error("Error while downloading video", "error", err, "path", path, "title", title, "video", video)
					continue
				}
				slog.Info("Downloaded video", "video", video, "path", path, "title", title)
				output <- DownlodedVideo{
					ScrapedVideo: video,
					SavePath:     path,
				}
			}
			close(output)
		}()
		return output
	}

	jobCreator := &JobCreator{
		Queries:  jm.Queries,
		Scrape:   scraper,
		VImport:  importer,
		Download: dowloader,
	}

	ticker := time.NewTicker(jm.AutoWakePeriod)

	for {
		select {
		case <-jm.Wake:
			jobCreator.RunJobPoolOnce()
		case <-ticker.C:
			jobCreator.RunJobPoolOnce()
		}
	}
}
