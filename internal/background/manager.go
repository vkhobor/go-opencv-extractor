package background

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/vkhobor/go-opencv/internal/config"
	"github.com/vkhobor/go-opencv/internal/db"
	"github.com/vkhobor/go-opencv/internal/mlog"
)

type DbMonitor struct {
	Wake        chan struct{}
	SqlDB       *sql.DB
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

	val, err := db.New(jm.SqlDB).GetVideosDownloadedButNotImported(context.Background())
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
