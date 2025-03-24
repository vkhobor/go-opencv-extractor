package features

import (
	"time"

	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/mlog"
	"github.com/vkhobor/go-opencv/queries"
	"github.com/vkhobor/go-opencv/youtube"
)

type DownloadVideoFeature struct {
	Queries *queries.Queries
	Config  config.DirectoryConfig
}

func (i *DownloadVideoFeature) DownloadVideo(videoID youtube.YoutubeVideo, jobID string) (savePath string, error error) {
	// TODO check if video is already downloaded, optionally abort while progressing
	mlog.Log().Info("Download started", "video", videoID, "job", jobID)

	savePath, err := i.retryDownloadWithAllClients(videoID)

	if err != nil {
		mlog.Log().Error("Error while downloading video", "error", err, "video", videoID)
		return savePath, err
	} else {
		mlog.Log().Info("Downloaded video", "video", videoID, "downloaded", savePath)
	}
	mlog.Log().Info("Download succesful", "video", videoID)

	err = i.Queries.SaveDownloadAttempt(string(videoID), savePath, err)
	if err != nil {
		mlog.Log().Error("Error while saving download attempt", "error", err, "videoID", videoID)
		return savePath, err
	}
	return
}

func (d *DownloadVideoFeature) handleProgress(progress chan float64, videoID youtube.YoutubeVideo) {
	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()

	for item := range progress {
		select {
		case <-ticker.C:
			mlog.Log().Info("downloadProgress", "id", videoID, "progress", item)
		default:
		}
	}
}

func (d *DownloadVideoFeature) retryDownloadWithAllClients(video youtube.YoutubeVideo) (string, error) {
	progress := make(chan float64)
	defer close(progress)
	go d.handleProgress(progress, video)

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
