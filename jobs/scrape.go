package jobs

import (
	"log/slog"
	"time"

	"github.com/vkhobor/go-opencv/domain"
	"github.com/vkhobor/go-opencv/scraper"
	"github.com/vkhobor/go-opencv/youtube"
)

type Scraper func()

// new Scraper creates a new Scraper
func NewScraper(queries *domain.JobQueries, throttle time.Duration) Scraper {
	scraper := scraper.Scraper{
		Throttle: throttle,
		Domain:   "yewtu.be",
	}
	scraperFunc := func(args domain.ScrapeArgs) []youtube.YoutubeVideo {
		return scraper.Scrape(args.Limit, args.SearchQuery)
	}

	return func() {
		toScrape := queries.GetToScrapeVideos()
		if len(toScrape) == 0 {
			slog.Debug("No videos to scrape")
			return
		}

		slog.Info("Running scrape job", "needed_to_scrape", toScrape)
		for _, scrapeArgs := range toScrape {
			toFind := scrapeArgs.Limit

			if toFind <= 0 {
				continue
			}

			scrapeChan := scraperFunc(domain.ScrapeArgs{SearchQuery: scrapeArgs.SearchQuery, Limit: scrapeArgs.Limit})

			for _, item := range scrapeChan {
				queries.SaveSraped(domain.ScrapedVideo{
					ID: item.String(),
				}, scrapeArgs.JobId)
			}
		}
	}

}
