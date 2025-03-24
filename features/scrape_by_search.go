package features

import (
	"errors"
	"fmt"
	"iter"

	"github.com/vkhobor/go-opencv/mlog"
	"github.com/vkhobor/go-opencv/queries"
	"github.com/vkhobor/go-opencv/scraper"
	"github.com/vkhobor/go-opencv/youtube"
)

type ScrapeBySearchFeature struct {
	scraper.Scraper
	Queries              *queries.Queries
	MaxErrorStopRetrying int
}

func (i *ScrapeBySearchFeature) ScrapeBySearch(keyword string, jobID string, limit int) (results []youtube.YoutubeVideo, error error) {

	// TODO check if video is already scraped, optionally abort while progressing

	results = []youtube.YoutubeVideo{}
	iterable := i.Scraper.AllForQuery(keyword)
	results, errs := i.collectScraped(iterable, limit)

	if len(errs) != 0 {
		mlog.Log().Error("Encountered errors while scraping", "errors", errs)
	}
	if len(results) < limit {
		return results, fmt.Errorf("could not scrape any videos")
	}

	// TODO make this more efficient, no need to query for every video
	jobs := i.Queries.GetToScrapeVideos()
	actualJob := queries.Job{}
	for _, job := range jobs {
		if job.JobID == jobID {
			actualJob = job
			break
		}
	}

	for _, result := range results {
		err := i.Queries.SaveNewlyScraped(jobID, string(result), actualJob.FilterID)
		if err != nil {
			if errors.Is(err, queries.ErrLimitExceeded) {
				mlog.Log().Warn("Reached scraping limit, stopping", "job", actualJob, "id", result)
				break
			} else if errors.Is(err, queries.ErrAlreadyScrapedForFilter) {
				mlog.Log().Warn("Collision, already scraped", "job", actualJob, "id", result)
				continue
			} else {
				return results, err
			}
		}
	}

	return
}

func (d *ScrapeBySearchFeature) collectScraped(
	scrape iter.Seq[scraper.Result],
	limit int) (collected []youtube.YoutubeVideo, errorsEncountered []error) {

	errorsEncountered = []error{}
	collected = []youtube.YoutubeVideo{}

	for result := range scrape {
		if len(collected) >= limit {
			return
		}
		if len(errorsEncountered) >= d.MaxErrorStopRetrying {
			break
		}

		if result.Error != nil {
			mlog.Log().Error("Error while scraping", "error", result.Error)
			errorsEncountered = append(errorsEncountered, result.Error)
			continue
		}

		mlog.Log().Debug("Scraped", "scraped", result)

		collected = append(collected, result.YoutubeVideo)
	}
	return
}
