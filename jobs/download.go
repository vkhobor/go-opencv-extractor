package jobs

import (
	"log/slog"
	"time"

	"github.com/vkhobor/go-opencv/domain"
	"github.com/vkhobor/go-opencv/youtube"
)

type Downloader func()

// new Downloader creates a new Downloader
func NewDownloader(queries *domain.JobQueries, throttle time.Duration) Downloader {

	return func() {
		slog.Info("Running download job")

		scraped := queries.GetScrapedVideos()
		if len(scraped) == 0 {
			slog.Debug("No videos to download")
			return
		}

		for _, video := range scraped {
			time.Sleep(throttle)
			slog.Info("Download started", "video", video)
			youtubeVideo := youtube.NewYoutubeIdFromScrapedVideo(video)
			downloaded := youtubeVideo.DownloadToFolder("videos")

			slog.Info("Downloaded video", "video", video, "downloaded", downloaded)
			queries.DownloadSaved(downloaded)
		}
	}
}
