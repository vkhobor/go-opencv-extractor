package scraper

import (
	"context"
	"database/sql"

	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/db"
)

type Queries struct {
	Queries *db.Queries
}

func (jc *Queries) GetToScrapeVideos() []ScrapeArgs {
	dbVal, err := jc.Queries.GetToScrapeVideos(context.Background())

	if err != nil {
		return []ScrapeArgs{}
	}

	return lo.FilterMap(dbVal, func(item db.GetToScrapeVideosRow, i int) (ScrapeArgs, bool) {
		return ScrapeArgs{
			SearchQuery: item.SearchQuery.String,
			Limit:       int(item.Limit.Int64 - item.FoundVideos),
			JobId:       item.ID,
		}, item.Limit.Int64-item.FoundVideos > 0
	})
}

func (jc *Queries) GetScrapedVideos() []ScrapedVideo {
	val, err := jc.Queries.GetScrapedVideos(context.Background())
	if err != nil {
		return []ScrapedVideo{}
	}

	result := make([]ScrapedVideo, len(val))
	for i, v := range val {
		result[i] = ScrapedVideo{ID: v.ID}
	}

	return result
}

func (jc *Queries) SaveSraped(video ScrapedVideo, jobId string) bool {
	_, err := jc.Queries.AddYtVideo(context.Background(), db.AddYtVideoParams{
		ID: video.ID,
		JobID: sql.NullString{
			String: jobId,
			Valid:  true,
		},
		Status: sql.NullString{
			String: "scraped",
			Valid:  true,
		},
	})

	if err != nil {
		return false
	}
	return true
}
