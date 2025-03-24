package background

import (
	"github.com/vkhobor/go-opencv/features"
	"github.com/vkhobor/go-opencv/mlog"
	"github.com/vkhobor/go-opencv/queries"
	"github.com/vkhobor/go-opencv/youtube"
)

func (d *DbMonitor) StartScrape() {
	for video := range d.ScrapeInput {
		if video.YouTubeID == "" {
			scraper := features.ScrapeByIdFeature{
				Queries:              d.Queries,
				MaxErrorStopRetrying: d.MaxErrorStopRetrying,
			}
			err := scraper.ScrapeByIdSearch(youtube.YoutubeVideo(video.YouTubeID), video.JobID, video.Limit)

			if err != nil {
				mlog.Log().Error("Error while scraping by id", "error", err)
				continue
			}
			d.DownloadInput <- queries.ScrapedVideo{
				Job: video,
				ID:  string(video.YouTubeID),
			}
		} else {
			scraper := features.ScrapeBySearchFeature{
				Scraper:              d.Scraper,
				Queries:              d.Queries,
				MaxErrorStopRetrying: d.MaxErrorStopRetrying,
			}
			results, err := scraper.ScrapeBySearch(video.SearchQuery, video.JobID, video.Limit)

			if err != nil {
				mlog.Log().Error("Error while scraping by search", "error", err)
				continue
			}

			for _, result := range results {
				d.DownloadInput <- queries.ScrapedVideo{
					Job: video,
					ID:  string(result),
				}
			}
		}
	}
}
