package background

import (
	"github.com/vkhobor/go-opencv/features"
	"github.com/vkhobor/go-opencv/mlog"
	"github.com/vkhobor/go-opencv/queries"
	"github.com/vkhobor/go-opencv/youtube"
)

func (d *DbMonitor) StartDownload() {
	for video := range d.DownloadInput {
		downloader := features.DownloadVideoFeature{
			Queries: d.Queries,
			Config:  d.Config,
		}

		savePath, err := downloader.DownloadVideo(youtube.YoutubeVideo(video.YouTubeID), video.JobID)
		if err != nil {
			mlog.Log().Error("Error while downloading video", "error", err, "video", video)
			continue
		}

		d.ImportInput <- queries.DownlodedVideo{
			ScrapedVideo: video,
			SavePath:     savePath,
			Error:        nil,
		}
	}
}
