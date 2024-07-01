package background

import (
	"log/slog"
	"sync"

	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/queries"
)

type DbMonitor struct {
	Wake            chan struct{}
	Queries         *queries.Queries
	ScrapeInput     chan<- queries.Job
	DownloadInput   chan<- queries.ScrapedVideo
	ImportInput     chan<- queries.DownlodedVideo
	cacheDownloaded BoundedQueue[queries.DownlodedVideo]
	cacheScraped    BoundedQueue[queries.ScrapedVideo]
	cacheJobs       BoundedQueue[queries.Job]
	once            sync.Once
}

func (jm *DbMonitor) Start() {
	jm.once.Do(func() {
		jm.cacheDownloaded = NewBoundedQueue[queries.DownlodedVideo](100)
		jm.cacheScraped = NewBoundedQueue[queries.ScrapedVideo](100)
		jm.cacheJobs = NewBoundedQueue[queries.Job](100)

		for {
			select {
			case <-jm.Wake:
				jm.PullWorkItemsFromDb()
			}
		}
	})
}

func (jm *DbMonitor) PullWorkItemsFromDb() {
	var waitGroup sync.WaitGroup
	waitGroup.Add(3)

	go func() {
		defer waitGroup.Done()

		jobs := jm.Queries.GetToScrapeVideos()
		jobs = lo.Filter(jobs, func(args queries.Job, index int) bool {
			return !jm.cacheJobs.Some(args)
		})

		slog.Debug("PullWorkItemsFromDb", "jobs", jobs)
		for _, args := range jobs {
			jm.cacheJobs.Push(args)
			jm.ScrapeInput <- args
		}
	}()

	go func() {
		defer waitGroup.Done()

		scrapedVideos := jm.Queries.GetScrapedVideos()
		scrapedVideos = lo.Filter(scrapedVideos, func(args queries.ScrapedVideo, index int) bool {
			return !jm.cacheScraped.Some(args)
		})

		slog.Debug("PullWorkItemsFromDb", "scrapedVideos", scrapedVideos)
		for _, video := range scrapedVideos {
			jm.cacheScraped.Push(video)
			jm.DownloadInput <- video
		}
	}()

	go func() {
		defer waitGroup.Done()

		downloadedVideos := jm.Queries.GetDownloadedVideos()
		downloadedVideos = lo.Filter(downloadedVideos, func(args queries.DownlodedVideo, index int) bool {
			return !jm.cacheDownloaded.Some(args)
		})

		slog.Debug("PullWorkItemsFromDb", "downloadedVideos", downloadedVideos)
		for _, video := range downloadedVideos {
			jm.cacheDownloaded.Push(video)
			jm.ImportInput <- video
		}
	}()

	waitGroup.Wait()
}
