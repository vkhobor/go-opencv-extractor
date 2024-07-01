package queries

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/vkhobor/go-opencv/db"
)

func (jc *Queries) GetDownloadedVideos() []DownlodedVideo {
	val, err := jc.Queries.GetVideosDownloaded(context.Background())
	if err != nil {
		slog.Error("GetDownloadedVideos: Error while getting downloaded videos", "error", err)
		return []DownlodedVideo{}
	}
	result := make([]DownlodedVideo, len(val))
	for i, v := range val {
		result[i] = DownlodedVideo{
			ScrapedVideo: ScrapedVideo{
				ID: v.YtVideoID,
				Job: Job{
					JobID:       v.JobID,
					SearchQuery: v.SearchQuery.String,
					FilterID:    v.FilterID.String,
				}},
			SavePath: v.Path,
		}
	}
	return result
}

func (jc *Queries) SaveDownloadAttempt(video DownlodedVideo) error {
	if video.Error != nil {
		err := jc.Queries.AddDownloadAttempt(context.Background(), db.AddDownloadAttemptParams{
			ID: uuid.New().String(),
			YtVideoID: sql.NullString{
				String: video.ID,
				Valid:  true,
			},
			Error: sql.NullString{
				String: video.Error.Error(),
				Valid:  true,
			},
			BlobStorageID: sql.NullString{
				Valid: false,
			},
		})
		return err
	}

	// TODO transaction
	blobId := uuid.New()
	err := jc.Queries.AddBlob(context.Background(), db.AddBlobParams{
		ID:   blobId.String(),
		Path: video.SavePath,
	})
	if err != nil {
		return err
	}

	err = jc.Queries.AddDownloadAttempt(context.Background(), db.AddDownloadAttemptParams{
		ID: uuid.New().String(),
		YtVideoID: sql.NullString{
			String: video.ID,
			Valid:  true,
		},
		Error: sql.NullString{
			Valid: false,
		},
		BlobStorageID: sql.NullString{
			String: blobId.String(),
			Valid:  true,
		},
	})
	return err
}

func RemoveAllPaths(files ...string) {
	for _, file := range files {
		_ = os.Remove(file)
	}
}
