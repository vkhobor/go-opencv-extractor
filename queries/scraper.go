package queries

import (
	"context"
	"database/sql"
	"errors"

	"github.com/vkhobor/go-opencv/db"
)

var ErrLimitExceeded = errors.New("over limit")
var ErrAlreadyScrapedForFilter = errors.New("already scraped for filter")

func (jc *Queries) SaveNewlyScraped(jobId string, videoID string, filterID string, name string) error {
	videoFromDb, err := jc.Queries.GetYtVideoWithJob(context.Background(), videoID)
	if err == nil && videoFromDb.FilterID.String == filterID {
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
		ID: videoID,
		Name: sql.NullString{
			String: name,
			Valid:  true,
		},
		JobID: sql.NullString{
			String: jobId,
			Valid:  true,
		},
	})

	return err
}
