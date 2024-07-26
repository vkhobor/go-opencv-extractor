package background

import (
	"errors"
	"fmt"
	"iter"

	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/mlog"
	"github.com/vkhobor/go-opencv/queries"
	"github.com/vkhobor/go-opencv/scraper"
)

type ScraperJob struct {
	scraper.Scraper
	Queries              *queries.Queries
	MaxErrorStopRetrying int
	Input                <-chan queries.Job
	Output               chan<- queries.ScrapedVideo
	Config               config.DirectoryConfig
}

func (d *ScraperJob) Start() {
	for video := range d.Input {
		results, err := d.scrapeSingle(video)
		if err != nil {
			mlog.Log().Error("Scraping did not produce any results", "error", err)
			continue
		}

		for _, result := range results {
			d.Output <- result
		}
	}
}

// TODO optionally move the single processing to another package e.g scrape/service
func (d *ScraperJob) scrapeSingle(args queries.Job) ([]queries.ScrapedVideo, error) {
	iterable := d.Scraper.AllForQuery(args.SearchQuery)
	results, errs := d.collectScraped(iterable, args)

	if len(errs) != 0 {
		mlog.Log().Error("Encountered errors while scraping", "errors", errs)
	}
	if len(results) < args.Limit {
		return nil, fmt.Errorf("could not scrape any videos", errors.Join(errs...))
	}

	return results, nil
}

func (d *ScraperJob) collectScraped(scrape iter.Seq[scraper.Result], args queries.Job) ([]queries.ScrapedVideo, []error) {
	errorsEncountered := []error{}
	collected := []queries.ScrapedVideo{}

	for result := range scrape {
		if len(errorsEncountered) >= d.MaxErrorStopRetrying {
			break
		}

		if result.Error != nil {
			mlog.Log().Error("Error while scraping", "error", result.Error)
			errorsEncountered = append(errorsEncountered, result.Error)
			continue
		}

		mlog.Log().Debug("Scraped", "scraped", result)
		scraped := queries.ScrapedVideo{
			Job: args,
			ID:  string(result.YoutubeVideo),
		}

		err := d.Queries.SaveNewlyScraped(scraped, args.JobID)
		if err != nil {
			if errors.Is(err, queries.ErrLimitExceeded) {
				mlog.Log().Warn("Reached scraping limit, stopping", "job", args, "id", scraped.ID)
				break
			} else if errors.Is(err, queries.ErrAlreadyScrapedForFilter) {
				mlog.Log().Warn("Collision, already scraped", "job", args, "id", scraped.ID)
				continue
			} else {
				// Unknown error
				errorsEncountered = append(errorsEncountered, err)
				continue
			}
		}

		collected = append(collected, scraped)
	}
	return collected, errorsEncountered
}
