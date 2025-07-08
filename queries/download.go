package queries

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/db"
)

func (jc *Queries) GetDownloadedVideos(includeImported bool) []DownlodedVideo {
	var result []DownlodedVideo
	if includeImported {
		val, err := jc.Queries.GetVideosDownloaded(context.Background())
		if err != nil {
			slog.Error("GetDownloadedVideos: Error while getting downloaded videos", "error", err)
			return []DownlodedVideo{}
		}
		result = make([]DownlodedVideo, len(val))
		for i, v := range val {
			result[i] = DownlodedVideo{
				ScrapedVideo: ScrapedVideo{
					ID: v.YtVideoID,
					Job: Job{
						JobID:       v.JobID,
						SearchQuery: v.SearchQuery.String,
						FilterID:    v.FilterID.String,
						Limit:       int(v.Limit.Int64),
						YouTubeID:   v.YtVideoID,
						Name:        v.YtVideoName.String,
					}},
				SavePath:       v.Path,
				ImportProgress: v.ImportProgress,
			}
		}
	} else {
		val, err := jc.Queries.GetVideosDownloadedButNotImported(context.Background())
		if err != nil {
			slog.Error("GetDownloadedVideos: Error while getting downloaded videos", "error", err)
			return []DownlodedVideo{}
		}
		result = make([]DownlodedVideo, len(val))
		for i, v := range val {
			result[i] = DownlodedVideo{
				ScrapedVideo: ScrapedVideo{
					ID: v.YtVideoID,
					Job: Job{
						JobID:       v.JobID,
						SearchQuery: v.SearchQuery.String,
						FilterID:    v.FilterID.String,
						Limit:       int(v.Limit.Int64),
						YouTubeID:   v.YtVideoID,
					}},
				SavePath:       v.Path,
				ImportProgress: v.ImportProgress.Int64,
			}
		}
	}
	return result
}

var ErrHasDownloaded = errors.New("already downloaded")

func (jc *Queries) SaveDownloadAttempt(videoID string, savePath string, downloadError error) error {
	attempts, err := jc.Queries.GetVideoWithDownloadAttempts(context.Background(), videoID)
	if err != nil {
		return err
	}

	if len(attempts) > 0 {
		hasSuccessful := lo.SomeBy(attempts, func(item db.GetVideoWithDownloadAttemptsRow) bool {
			return !item.Error.Valid
		})
		if hasSuccessful {
			return ErrHasDownloaded
		}
	}

	if downloadError != nil {
		err := jc.Queries.AddDownloadAttempt(context.Background(), db.AddDownloadAttemptParams{
			ID: uuid.New().String(),
			YtVideoID: sql.NullString{
				String: videoID,
				Valid:  true,
			},
			Error: sql.NullString{
				String: downloadError.Error(),
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
	err = jc.Queries.AddBlob(context.Background(), db.AddBlobParams{
		ID:   blobId.String(),
		Path: savePath,
	})
	if err != nil {
		return err
	}

	err = jc.Queries.AddDownloadAttempt(context.Background(), db.AddDownloadAttemptParams{
		ID: uuid.New().String(),
		YtVideoID: sql.NullString{
			String: videoID,
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
