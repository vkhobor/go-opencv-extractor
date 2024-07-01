package background

import (
	"errors"
	"log/slog"

	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/queries"
	"github.com/vkhobor/go-opencv/scraper"
	"github.com/vkhobor/go-opencv/youtube"
)

type ScraperJob struct {
	scraper.Scraper
	Queries *queries.Queries
	Input   <-chan queries.Job
	Output  chan<- queries.ScrapedVideo
	Config  config.DirectoryConfig
}

// Start starts the Scraper
func (d *ScraperJob) Start() {
	for video := range d.Input {
		results, err := d.scrapeSingle(video)
		if err != nil {
			slog.Error(
				"Error while scraping video",
				"error", err,
				"video", video,
				"method", "Start")
		}

		for _, result := range results {
			d.Output <- result
		}
	}
}

func (d *ScraperJob) scrapeSingle(args queries.Job) ([]queries.ScrapedVideo, error) {
	if args.Limit <= 0 {
		return nil, errors.New("limit is less than or equal to 0")
	}
	limit := args.Limit
	if limit > 50 {
		limit = 50
	}

	results := []queries.ScrapedVideo{}
	err := d.Scraper.Scrape(args.SearchQuery, d.handleFound(args, limit, results))

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (d *ScraperJob) handleFound(args queries.Job, limit int, output []queries.ScrapedVideo) func(item youtube.YoutubeVideo, err error, stop func()) {
	saved := 0
	errored := 0

	f := func(item youtube.YoutubeVideo, err error, stop func()) {
		if errored >= 5 {
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

		// TODO should only allow saving if does not make the database inconsistent, like scraping over the limit
		ok := d.Queries.SaveNewlyScraped(scraped, args.JobID)
		if !ok {
			slog.Error("Error while saving scraped video", "video", scraped, "method", "scrapeSingle")
			errored++
		}

		saved++
		output = append(output, scraped)

		if saved >= limit {
			d.Scraper.Stop()
		}
	}
	return f
}

func (d *ScraperJob) handleAlreadySaved(video queries.ScrapedVideo) error {
	scraped := d.Queries.GetScrapedVideos()
	found, ok := lo.Find(scraped, func(item queries.ScrapedVideo) bool {
		return item.ID == video.ID
	})
	if !ok {
		return errors.New("video not found")
	}

	if found.Job.FilterID == video.FilterID {
		return nil
	}

	// TODO attach the job's filter request to the video

	slog.Debug("Already saved", "video", video, "found", found, "method", "handleAlreadySaved")
	return nil
}
