package scraper

import (
	"context"
	"errors"
	"log/slog"

	"github.com/vkhobor/go-opencv/config"
)

type Job struct {
	Limit       int
	JobID       string
	SearchQuery string
	FilterID    string
}

type ScrapedVideo struct {
	Job
	ID string
}

type ScraperJob struct {
	Scraper
	Queries *Queries
	Input   <-chan Job
	Output  chan<- ScrapedVideo
	Config  config.DirectoryConfig
}

// Start starts the Scraper
func (d *ScraperJob) Start() {
	for video := range d.Input {
		err := d.scrapeSingle(video)
		if err != nil {
			slog.Error("Error while importing video", "error", err, "video", video)
		}
	}
}

func (d *ScraperJob) scrapeSingle(args Job) error {
	if args.Limit <= 0 {
		return errors.New("limit is less than or equal to 0")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// This is chan so later we can test not to save already saved videos
	scrapeChan, err := d.Scraper.Scrape(ctx, args.SearchQuery)
	if err != nil {
		return err
	}

	saved := 0
	scraped := 0
	for item := range scrapeChan {
		slog.Debug("Scraped", "scraped", item)
		scraped++

		if saved >= args.Limit {
			break
		}

		if scraped >= 50 {
			break
		}

		// TODO if already saved, only attach the job's filter request if not attached already to video
		scraped := ScrapedVideo{
			Job: args,
			ID:  item.String(),
		}
		ok := d.Queries.SaveSraped(scraped, args.JobID)
		if ok {
			saved++
			d.Output <- scraped
		}
	}
	return nil
}
