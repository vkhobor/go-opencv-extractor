package background

import (
	"errors"
	"log/slog"

	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/queries"
	"github.com/vkhobor/go-opencv/scraper"
	"github.com/vkhobor/go-opencv/youtube"
)

type ScraperJob struct {
	scraper.Scraper
	Queries              *queries.Queries
	MaxErrorStopRetrying int
	Input                <-chan queries.Job
	Output               chan<- queries.ScrapedVideo
	Config               config.DirectoryConfig
	WakeJobs             chan<- struct{}
}

// Start starts the Scraper
func (d *ScraperJob) Start() {
	for video := range d.Input {
		results, err := d.scrapeSingle(video)
		if err != nil {
			continue
		}

		for _, result := range results {
			select {
			// Try to send the downloaded video to the output channel, potentially saving a database call
			case d.Output <- result:
			// If the output channel is full, try to wake the job manager
			case d.WakeJobs <- struct{}{}:
			// If the job manager is awake it will pick it up from the database on next pull from channel
			default:
			}
		}
	}
}

func (d *ScraperJob) scrapeSingle(args queries.Job) ([]queries.ScrapedVideo, error) {

	results := []queries.ScrapedVideo{}
	err := d.Scraper.Scrape(args.SearchQuery, d.handleFound(args, &results))

	if err != nil {
		slog.Error("Error while setting up scraper", "error", err, "method", "scrapeSingle")
		return nil, err
	}

	return results, nil
}

func (d *ScraperJob) handleFound(args queries.Job, output *[]queries.ScrapedVideo) func(item youtube.YoutubeVideo, err error, stop func()) {
	errored := 0

	return func(item youtube.YoutubeVideo, err error, stop func()) {
		if errored >= d.MaxErrorStopRetrying {
			stop()
			return
		}

		if err != nil {
			slog.Error("Error while scraping", "error", err)
			errored++
			return
		}

		slog.Debug("Scraped", "scraped", item)
		scraped := queries.ScrapedVideo{
			Job: args,
			ID:  item.String(),
		}

		err = d.Queries.SaveNewlyScraped(scraped, args.JobID)
		if err != nil {
			slog.Error("Error while saving scraped video", "error", err, "video", scraped, "method", "scrapeSingle")

			if errors.Is(err, queries.ErrLimitExceeded) {
				stop()
				return
			} else if errors.Is(err, queries.ErrAlreadyScrapedForFilter) {
				return
			} else {
				// Unknown error
				errored++
			}
		}

		*output = append(*output, scraped)
	}
}
