package background

import (
	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/mlog"
	"github.com/vkhobor/go-opencv/queries"
)

type DbMonitor struct {
	Wake        chan struct{}
	Queries     *queries.Queries
	ImportInput chan queries.DownlodedVideo
	Config      config.DirectoryConfig
}

func (jm *DbMonitor) Start() {
	go jm.StartImport()

	for range jm.Wake {
		jm.PullWorkItemsFromDb()
	}
}

func (jm *DbMonitor) PullWorkItemsFromDb() {
	downloadedVideos := jm.Queries.GetDownloadedVideos(false)

	mlog.Log().Debug("PullWorkItemsFromDb", "downloadedVideos", downloadedVideos)
	for _, video := range downloadedVideos {
		jm.ImportInput <- video
	}
}
