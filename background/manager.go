package background

import (
	"sync"

	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/mlog"
	"github.com/vkhobor/go-opencv/queries"
	"github.com/vkhobor/go-opencv/scraper"
)

type DbMonitor struct {
	Wake    chan struct{}
	Queries *queries.Queries
	scraper.Scraper
	ScrapeInput          chan queries.Job
	DownloadInput        chan queries.ScrapedVideo
	ImportInput          chan queries.DownlodedVideo
	Config               config.DirectoryConfig
	MaxErrorStopRetrying int
}

func (jm *DbMonitor) Start() {
	go jm.StartDownload()
	go jm.StartImport()
	go jm.StartScrape()

	for range jm.Wake {
		jm.PullWorkItemsFromDb()
	}
}

func (jm *DbMonitor) PullWorkItemsFromDb() {
	var waitGroup sync.WaitGroup
	waitGroup.Add(3)

	go func() {
		defer waitGroup.Done()

		jobs := jm.Queries.GetToScrapeVideos()

		mlog.Log().Debug("PullWorkItemsFromDb", "jobs", jobs)
		for _, args := range jobs {
			jm.ScrapeInput <- args
		}
	}()

	go func() {
		defer waitGroup.Done()

		scrapedVideos := jm.Queries.GetScrapedVideos()

		mlog.Log().Debug("PullWorkItemsFromDb", "scrapedVideos", scrapedVideos)
		for _, video := range scrapedVideos {
			jm.DownloadInput <- video
		}
	}()

	go func() {
		defer waitGroup.Done()

		downloadedVideos := jm.Queries.GetDownloadedVideos()

		mlog.Log().Debug("PullWorkItemsFromDb", "downloadedVideos", downloadedVideos)
		for _, video := range downloadedVideos {
			jm.ImportInput <- video
		}
	}()

	waitGroup.Wait()
}
