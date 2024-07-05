package queries

import (
	"context"
	"database/sql"
	"errors"

	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/db"
)

func (jc *Queries) GetToScrapeVideos() []Job {

	// Hi
	dbVal, err := jc.Queries.GetJobs(context.Background())

	if err != nil {
		return []Job{}
	}

	return lo.FilterMap(dbVal, func(item db.GetJobsRow, i int) (Job, bool) {
		return Job{
			FilterID:    item.FilterID.String,
			SearchQuery: item.SearchQuery.String,
			Limit:       int(item.Limit.Int64 - item.VideosFound),
			JobID:       item.ID,
		}, item.Limit.Int64-item.VideosFound > 0
	})
}

func (jc *Queries) GetScrapedVideos() []ScrapedVideo {
	val, err := jc.Queries.GetScrapedVideos(context.Background())
	if err != nil {
		return []ScrapedVideo{}
	}

	result := make([]ScrapedVideo, len(val))
	for i, v := range val {
		result[i] = ScrapedVideo{ID: v.YtVideoID}
	}

	return result
}

var ErrLimitExceeded = errors.New("over limit")
var ErrAlreadyScrapedForFilter = errors.New("already scraped for filter")

func (jc *Queries) SaveNewlyScraped(video ScrapedVideo, jobId string) error {
	videoFromDb, err := jc.Queries.GetYtVideoWithJob(context.Background(), video.ID)
	if err == nil && videoFromDb.FilterID.String == video.FilterID {
		return ErrAlreadyScrapedForFilter
	} else if err == nil {
		// TODO connect to filter or job if multiple filters can exist
	} else if err != nil && err != sql.ErrNoRows {
		return err
	}

	job, err := jc.Queries.GetJob(context.Background(), jobId)
	if err != nil {
		return err
	}

	if job.VideosFound >= job.Limit.Int64 {
		return ErrLimitExceeded
	}

	_, err = jc.Queries.AddYtVideo(context.Background(), db.AddYtVideoParams{
		ID: video.ID,
		JobID: sql.NullString{
			String: jobId,
			Valid:  true,
		},
	})

	return err
}
