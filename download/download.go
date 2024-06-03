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

	savePath, err := d.RetryDownloadWithAllClients(youtubeVideo, video)

	if err != nil {
		slog.Error("Error while downloading video", "error", err, "video", video)

		return DownlodedVideo{ScrapedVideo: video, Error: err}
	}

	downloaded := DownlodedVideo{
		ScrapedVideo: video,
		SavePath:     savePath,
		Error:        nil,
	}
	d.Queries.SaveDownloadAttempt(downloaded)

	slog.Info("Downloaded video", "video", video, "downloaded", downloaded)

	return downloaded
}

func HandleProgress(progress chan float64, video scraper.ScrapedVideo) {
	for p := range progress {
		slog.Info("Download progress", "video", video, "progress", p)
	}
}

func (d *Downloader) RetryDownloadWithAllClients(video youtube.YoutubeVideo, scrapedVideo scraper.ScrapedVideo) (string, error) {
	progress := make(chan float64)
	defer close(progress)
	go HandleProgress(progress, scrapedVideo)

	savePath, err := video.DownloadToFolder(youtube.AndroidClient, d.Config.GetVideosDir(), progress)
	if err == nil {
		return savePath, nil
	}
	

	savePath, err = video.DownloadToFolder(youtube.WebClient, d.Config.GetVideosDir(), progress)
	if err == nil {
		return savePath, nil
	}

	savePath, err = video.DownloadToFolder(youtube.EmbeddedClient, d.Config.GetVideosDir(), progress)
	if err == nil {
		return savePath, nil
	}

	return "", err
}
