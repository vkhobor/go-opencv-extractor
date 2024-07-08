package background

import (
	"log/slog"
	"sync"

	"github.com/vkhobor/go-opencv/queries"
)

type DbMonitor struct {
	// TODO separate into different channels for each type of work
	Wake          chan struct{}
	Queries       *queries.Queries
	ScrapeInput   chan<- queries.Job
	DownloadInput chan<- queries.ScrapedVideo
	ImportInput   chan<- queries.DownlodedVideo
}

func (jm *DbMonitor) Start() {
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

		slog.Debug("PullWorkItemsFromDb", "jobs", jobs)
		for _, args := range jobs {
			jm.ScrapeInput <- args
		}
	}()

	go func() {
		defer waitGroup.Done()

		scrapedVideos := jm.Queries.GetScrapedVideos()

		slog.Debug("PullWorkItemsFromDb", "scrapedVideos", scrapedVideos)
		for _, video := range scrapedVideos {
			jm.DownloadInput <- video
		}
	}()

	go func() {
		defer waitGroup.Done()

		downloadedVideos := jm.Queries.GetDownloadedVideos()

		slog.Debug("PullWorkItemsFromDb", "downloadedVideos", downloadedVideos)
		for _, video := range downloadedVideos {
			jm.ImportInput <- video
		}
	}()

	waitGroup.Wait()
}
