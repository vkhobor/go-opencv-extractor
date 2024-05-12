package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/vkhobor/go-opencv/domain"
	"github.com/vkhobor/go-opencv/scraper"
)

type Scraper func()

// new Scraper creates a new Scraper
func NewScraper(queries *domain.JobQueries, throttle time.Duration) Scraper {
	scraper := scraper.Scraper{
		Throttle: throttle,
		Domain:   "yewtu.be",
	}

	return func() {
		toScrape := queries.GetToScrapeVideos()
		if len(toScrape) == 0 {
			slog.Debug("No videos to scrape")
			return
		}

		slog.Info("Running scrape job", "needed_to_scrape", toScrape)
		for _, args := range toScrape {
			toFind := args.Limit

			if toFind <= 0 {
				continue
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// This is chan so later we can test not to save already saved videos
			scrapeChan, err := scraper.Scrape(ctx, args.SearchQuery)
			if err != nil {
				slog.Error("Error scraping", "error", err)
				continue
			}

			saved := 0
			scraped := 0
			for item := range scrapeChan {
				slog.Debug("Scraped", "scraped", item)
				scraped++

				if saved >= toFind {
					break
				}

				if scraped >= 50 {
					break
				}

				// TODO if already saved, only attach the job's filter request if not attached already to video
				ok := queries.SaveSraped(domain.ScrapedVideo{
					ID: item.String(),
				}, args.JobId)
				if ok {
					saved++
				}
			}
		}
	}

}
