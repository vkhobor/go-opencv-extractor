package download

import (
	"log/slog"
	"time"

	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/scraper"
	"github.com/vkhobor/go-opencv/youtube"
)

type DownlodedVideo struct {
	scraper.ScrapedVideo
	SavePath string
	Error    error
}

type Downloader struct {
	Queries  *Queries
	Throttle time.Duration
	Config   config.DirectoryConfig
	Input    <-chan scraper.ScrapedVideo
	Output   chan<- DownlodedVideo
}

func (d *Downloader) Start() {
	for video := range d.Input {
		d.Output <- d.downloadVideo(video)
	}
}

func (d *Downloader) downloadVideo(video scraper.ScrapedVideo) DownlodedVideo {
	time.Sleep(d.Throttle)
	slog.Info("Download started", "video", video)
	youtubeVideo := youtube.YoutubeVideo(video.ID)

	progress := make(chan float64)
	defer close(progress)
	go HandleProgress(progress, video)

	savePath, err := youtubeVideo.DownloadToFolder(d.Config.GetVideosDir(), progress)
	if err != nil {
		slog.Error("Error while downloading video", "error", err, "video", video)
		return DownlodedVideo{ScrapedVideo: video, Error: err}
	}

	downloaded := DownlodedVideo{
		ScrapedVideo: video,
		SavePath:     savePath,
		Error:        nil,
	}
	d.Queries.DownloadSaved(downloaded)

	slog.Info("Downloaded video", "video", video, "downloaded", downloaded)

	return downloaded
}

func HandleProgress(progress chan float64, video scraper.ScrapedVideo) {
	for p := range progress {
		slog.Info("Download progress", "video", video, "progress", p)
	}
}
