package background

import (
	"context"
	"log/slog"

	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/db"
	"github.com/vkhobor/go-opencv/mlog"
)

type DbMonitor struct {
	Wake        chan struct{}
	Queries     *db.Queries
	ImportInput chan struct {
		ID       string
		JobID    string
		FilterID string
	}
	Config config.DirectoryConfig
}

func (jm *DbMonitor) Start() {
	go jm.StartImport()

	for range jm.Wake {
		jm.PullWorkItemsFromDb()
	}
}

func (jm *DbMonitor) PullWorkItemsFromDb() {
	val, err := jm.Queries.GetVideosDownloadedButNotImported(context.Background())
	if err != nil {
		slog.Error("GetDownloadedVideos: Error while getting downloaded videos", "error", err)
		return
	}

	mlog.Log().Debug("PullWorkItemsFromDb", "downloadedVideos", val)
	for _, video := range val {
		jm.ImportInput <- struct {
			ID       string
			JobID    string
			FilterID string
		}{
			ID:       video.YtVideoID,
			JobID:    video.JobID,
			FilterID: video.FilterID.String,
		}
	}
}
