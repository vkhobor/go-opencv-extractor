package domain

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/db"
)

type JobQueries struct {
	Queries *db.Queries
}

func (jc *JobQueries) GetToScrapeVideos() []ScrapeArgs {
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

func (jc *JobQueries) GetScrapedVideos() []ScrapedVideo {
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

func (jc *JobQueries) GetDownloadedVideos() []DownlodedVideo {
	val, err := jc.Queries.GetVideosDownloaded(context.Background())
	if err != nil {
		return []DownlodedVideo{}
	}
	result := make([]DownlodedVideo, len(val))
	for i, v := range val {
		result[i] = DownlodedVideo{ScrapedVideo: ScrapedVideo{ID: v.ID}, SavePath: v.Path}
	}
	return result
}

func (jc *JobQueries) GetRefImages() ([]string, error) {
	val, err := jc.Queries.GetReferences(context.Background())
	if err != nil {
		return nil, err
	}

	return lo.Map(val, func(item db.GetReferencesRow, i int) string {
		return item.Path
	}), nil
}

func (jc *JobQueries) SaveSraped(video ScrapedVideo, jobId string) bool {
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

func (jc *JobQueries) DownloadSaved(video DownlodedVideo) {
	slog.Debug("Saving downloaded", "video", video)
	if video.Error != nil {
		_, err := jc.Queries.UpdateStatus(context.Background(), db.UpdateStatusParams{
			ID: video.ID,
			Status: sql.NullString{
				String: "errored",
				Valid:  true,
			},
			Error: sql.NullString{
				Valid:  true,
				String: video.Error.Error(),
			},
		})
		if err != nil {
			slog.Error("Error while updating status", "error", err)
		}

		return
	}
	blobId := uuid.New()
	_, errAddBlob := jc.Queries.AddBlob(context.Background(), db.AddBlobParams{
		ID:   blobId.String(),
		Path: video.SavePath,
	})

	_, errUpdateStatus := jc.Queries.UpdateStatus(context.Background(), db.UpdateStatusParams{
		ID: video.ID,
		Status: sql.NullString{
			String: "downloaded",
			Valid:  true,
		},
	})

	jc.Queries.AddBlobToVideo(context.Background(), db.AddBlobToVideoParams{
		BlobStorageID: sql.NullString{
			String: blobId.String(),
			Valid:  true,
		},
		ID: video.ID,
	})

	if errAddBlob != nil || errUpdateStatus != nil {
		slog.Error("Error while updating status", "error", errAddBlob, "error2", errUpdateStatus)
		RemoveAllPaths(video.SavePath)
		return
	}
}

func (jc *JobQueries) SaveImported(video ImportedVideo) {
	if video.Error != nil {
		jc.Queries.UpdateStatus(context.Background(), db.UpdateStatusParams{
			ID: video.ID,
			Status: sql.NullString{
				String: "errored",
				Valid:  true,
			},
		})
		return
	}

	paths := []string{}
	for _, frame := range video.ExtractedFrames {
		paths = append(paths, frame.Path)

		blobID := uuid.New()
		jc.Queries.AddBlob(context.Background(), db.AddBlobParams{
			ID:   blobID.String(),
			Path: frame.Path,
		})
		jc.Queries.AddPicture(context.Background(), db.AddPictureParams{
			ID: uuid.New().String(),
			YtVideoID: sql.NullString{
				String: video.ID,
				Valid:  true,
			},
			FrameNumber: sql.NullInt64{
				Int64: int64(frame.FrameNumber),
				Valid: true,
			},
			BlobStorageID: sql.NullString{
				String: blobID.String(),
				Valid:  true,
			},
		})
	}

	jc.Queries.UpdateStatus(
		context.Background(),
		db.UpdateStatusParams{
			Status: sql.NullString{
				String: "imported",
				Valid:  true,
			},
			ID: video.ID,
		})

	return
}

func RemoveAllPaths(files ...string) {
	for _, file := range files {
		_ = os.Remove(file)
	}
}
