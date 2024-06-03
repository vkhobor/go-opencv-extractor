package jobs

import (
	"time"

	"github.com/vkhobor/go-opencv/download"
	"github.com/vkhobor/go-opencv/scraper"
)

type DbMonitor struct {
	Wake            chan struct{}
	AutoWakePeriod  time.Duration
	ScrapeQueries   *scraper.Queries
	DownloadQueries *download.Queries
	ScrapeInput     chan<- scraper.Job
	DownloadInput   chan<- scraper.ScrapedVideo
	ImportInput     chan<- download.DownlodedVideo
}

func (jm *DbMonitor) Start() {

	ticker := time.NewTicker(jm.AutoWakePeriod)

	for {
		select {
		case <-jm.Wake:
			jm.PullWorkItemsFromDb()
		case <-ticker.C:
			jm.PullWorkItemsFromDb()
		}
	}
}

func (jm *DbMonitor) PullWorkItemsFromDb() {
	go func() {
		scrapeArgs := jm.ScrapeQueries.GetToScrapeVideos()
		for _, args := range scrapeArgs {
			jm.ScrapeInput <- args
		}
	}()

	go func() {
		scrapedVideos := jm.ScrapeQueries.GetScrapedVideos()
		for _, video := range scrapedVideos {
			jm.DownloadInput <- video
		}
	}()

	go func() {
		downloadedVideos := jm.DownloadQueries.GetDownloadedVideos()
		for _, video := range downloadedVideos {
			jm.ImportInput <- video
		}
	}()
}
