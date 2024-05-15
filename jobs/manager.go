package jobs

import (
	"time"

	"github.com/vkhobor/go-opencv/download"
	"github.com/vkhobor/go-opencv/scraper"
)

type JobManager struct {
	Wake            chan struct{}
	AutoWakePeriod  time.Duration
	ScrapeQueries   *scraper.Queries
	DownloadQueries *download.Queries
	ScrapeInput     chan<- scraper.ScrapeArgs
	DownloadInput   chan<- scraper.ScrapedVideo
	ImportInput     chan<- download.DownlodedVideo
}

func (jm *JobManager) Start() {

	ticker := time.NewTicker(jm.AutoWakePeriod)

	for {
		select {
		case <-jm.Wake:
			jm.RunPipelineOnce()
		case <-ticker.C:
			jm.RunPipelineOnce()
		}
	}
}

func (jm *JobManager) RunPipelineOnce() {
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
