package queries

import (
	"context"
	"database/sql"
	"errors"

	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/db"
)

func (jc *Queries) GetToScrapeVideos() []Job {
	dbVal, err := jc.Queries.GetJobs(context.Background())

	if err != nil {
		return []Job{}
	}

	return lo.FilterMap(dbVal, func(item db.GetJobsRow, i int) (Job, bool) {
		return Job{
			FilterID:    item.FilterID.String,
			SearchQuery: item.SearchQuery.String,
			Limit:       int(item.Limit.Int64 - item.FoundVideos),
			JobID:       item.ID,
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
		result[i] = ScrapedVideo{ID: v.YtVideoID}
	}

	return result
}

var ErrLimitExceeded = errors.New("over limit")

func (jc *Queries) SaveNewlyScraped(video ScrapedVideo, jobId string) error {
	job, err := jc.Queries.GetJob(context.Background(), jobId)
	if err != nil {
		return err
	}

	if job.VideosFound >= job.Limit.Int64 {
		return ErrLimitExceeded
	}

	// TODO if vieo already exists, update it, connect it to the job and return
	_, err = jc.Queries.AddYtVideo(context.Background(), db.AddYtVideoParams{
		ID: video.ID,
		JobID: sql.NullString{
			String: jobId,
			Valid:  true,
		},
	})

	return err
}
