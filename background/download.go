package background

import (
	"time"

	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/mlog"
	"github.com/vkhobor/go-opencv/queries"
	"github.com/vkhobor/go-opencv/youtube"
)

type Downloader struct {
	Queries  *queries.Queries
	Throttle time.Duration
	Config   config.DirectoryConfig
	Input    <-chan queries.ScrapedVideo
	Output   chan<- queries.DownlodedVideo
}

func (d *Downloader) Start() {
	for video := range d.Input {
		downloaded, err := d.downloadVideo(video)
		if err != nil {
			mlog.Log().Error("Error while downloading video", "error", err, "video", video)
			continue
		}

		d.Output <- downloaded
	}
}

// TODO optionally move the single processing to another package e.g download/service
func (d *Downloader) downloadVideo(video queries.ScrapedVideo) (queries.DownlodedVideo, error) {
	// TODO check if video is already downloaded, optionally abort while progressing
	time.Sleep(d.Throttle)
	mlog.Log().Info("Download started", "video", video)
	youtubeVideo := youtube.YoutubeVideo(video.ID)

	savePath, err := d.RetryDownloadWithAllClients(youtubeVideo, video)
	downloaded := queries.DownlodedVideo{
		ScrapedVideo: video,
		SavePath:     savePath,
		Error:        err,
	}

	if err != nil {
		mlog.Log().Error("Error while downloading video", "error", err, "video", video)
		return queries.DownlodedVideo{}, err
	} else {
		mlog.Log().Info("Downloaded video", "video", video, "downloaded", downloaded)
	}
	mlog.Log().Info("Download succesful", "video", video)

	err = d.Queries.SaveDownloadAttempt(downloaded)
	if err != nil {
		mlog.Log().Error("Error while saving download attempt", "error", err, "video", video)
		return queries.DownlodedVideo{}, err
	}

	return downloaded, nil
}

func (d *Downloader) HandleProgress(progress chan float64, video queries.ScrapedVideo) {
	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()

	for item := range progress {
		select {
		case <-ticker.C:
			mlog.Log().Info("downloadProgress", "id", video.ID, "progress", item)
		default:
		}
	}
}

func (d *Downloader) RetryDownloadWithAllClients(video youtube.YoutubeVideo, scrapedVideo queries.ScrapedVideo) (string, error) {
	progress := make(chan float64)
	defer close(progress)
	go d.HandleProgress(progress, scrapedVideo)

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
