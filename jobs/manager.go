package jobs

import (
	"log/slog"
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
		slog.Debug("PullWorkItemsFromDb", "scrapeArgs", scrapeArgs)
		for _, args := range scrapeArgs {
			jm.ScrapeInput <- args
		}
	}()

	go func() {
		scrapedVideos := jm.ScrapeQueries.GetScrapedVideos()
		slog.Debug("PullWorkItemsFromDb", "scrapedVideos", scrapedVideos)
		for _, video := range scrapedVideos {
			jm.DownloadInput <- video
		}
	}()

	go func() {
		downloadedVideos := jm.DownloadQueries.GetDownloadedVideos()
		slog.Debug("PullWorkItemsFromDb", "downloadedVideos", downloadedVideos)
		for _, video := range downloadedVideos {
			jm.ImportInput <- video
		}
	}()
}
