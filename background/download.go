package background

import (
	"log/slog"
	"time"

	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/queries"
	"github.com/vkhobor/go-opencv/youtube"
)

type Downloader struct {
	Queries  *queries.Queries
	Throttle time.Duration
	Config   config.DirectoryConfig
	Input    <-chan queries.ScrapedVideo
	Output   chan<- queries.DownlodedVideo
	WakeJobs chan<- struct{}
}

func (d *Downloader) Start() {
	for video := range d.Input {
		downloaded, err := d.downloadVideo(video)
		if err != nil {
			slog.Error("Error while downloading video", "error", err, "video", video)
			continue
		}

		select {
		case d.Output <- downloaded:
			continue
		default:
		}

		select {
		case d.WakeJobs <- struct{}{}:
			continue
		default:
		}
	}
}

func (d *Downloader) downloadVideo(video queries.ScrapedVideo) (queries.DownlodedVideo, error) {
	time.Sleep(d.Throttle)
	slog.Info("Download started", "video", video)
	youtubeVideo := youtube.YoutubeVideo(video.ID)

	savePath, err := d.RetryDownloadWithAllClients(youtubeVideo, video)
	downloaded := queries.DownlodedVideo{
		ScrapedVideo: video,
		SavePath:     savePath,
		Error:        err,
	}

	if err != nil {
		slog.Error("Error while downloading video", "error", err, "video", video)
		return queries.DownlodedVideo{}, err
	} else {
		slog.Info("Downloaded video", "video", video, "downloaded", downloaded)
	}

	err = d.Queries.SaveDownloadAttempt(downloaded)
	if err != nil {
		slog.Error("Error while saving download attempt", "error", err, "video", video)
		return queries.DownlodedVideo{}, err
	}

	return downloaded, nil
}

func HandleProgress(progress chan float64, video queries.ScrapedVideo) {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for item := range progress {
		select {
		case <-ticker.C:
			slog.Info("downloadProgress", "id", video.ID, "progress", item)
		default:
		}
	}
}

func (d *Downloader) RetryDownloadWithAllClients(video youtube.YoutubeVideo, scrapedVideo queries.ScrapedVideo) (string, error) {
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
