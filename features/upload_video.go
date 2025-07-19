package features

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/db"
	"github.com/vkhobor/go-opencv/mlog"
)

type UploadVideoFeature struct {
	DbSql    TXer
	Querier  QuerierWithTx
	Config   config.DirectoryConfig
	WakeJobs chan<- struct{}
}

func (i *UploadVideoFeature) DownloadVideo(ctx context.Context, data io.Reader, filterId string, name string) (savePath string, error error) {
	tx, err := i.DbSql.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	queries := i.Querier.WithTx(tx)

	jobID := uuid.New().String()
	mlog.Log().Info("Creating new job", "jobID", jobID, "filterID", filterId)
	_, err = queries.CreateJob(ctx, db.CreateJobParams{
		FilterID: sql.NullString{
			String: filterId,
			Valid:  true,
		},
		Limit: sql.NullInt64{
			Valid: true,
			Int64: 1,
		},
		ID: jobID,
	})

	videoId := uuid.New().String()
	mlog.Log().Info("Saving newly scraped video", "jobID", jobID, "videoID", videoId, "filterID", filterId)
	err = i.SaveNewlyScraped(ctx, tx, jobID, videoId, filterId, name)
	if err != nil {
		mlog.Log().Error("Failed to save newly scraped video", "error", err)
		return "", err
	}

	folderPath := i.Config.GetVideosDir()
	mlog.Log().Debug("Creating videos directory", "path", folderPath)
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		mlog.Log().Error("Failed to create videos directory", "error", err)
		return "", err
	}

	id := uuid.New()
	fileName := fmt.Sprintf("%v_%v.mp4", id.String(), videoId)
	filePath := filepath.Join(folderPath, fileName)
	mlog.Log().Info("Preparing to save video file", "path", filePath)

	savePath = filePath
	dst, err := os.Create(filePath)
	if err != nil {
		mlog.Log().Error("Failed to create video file", "error", err, "path", filePath)
		return "", err
	}
	defer dst.Close()

	mlog.Log().Info("Starting video data copy", "videoID", videoId)
	buffer := make([]byte, 1024)
	bytesWritten := 0
	for {
		n, err := data.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			mlog.Log().Error("Error reading video data", "error", err)
			return "", err
		}
		dst.Write(buffer[:n])
		bytesWritten += n
	}
	mlog.Log().Info("Completed video data copy", "videoID", videoId, "bytesWritten", bytesWritten)

	err = i.SaveDownloadAttempt(ctx, tx, videoId, filePath, err)
	if err != nil {
		mlog.Log().Error("Error while saving download attempt", "error", err, "videoID", videoId)
		return savePath, err
	}

	mlog.Log().Info("Successfully completed video download", "videoID", videoId, "path", filePath)

	select {
	case i.WakeJobs <- struct{}{}:
		mlog.Log().Info("Waking up jobs")
	default:
		mlog.Log().Info("Jobs already awake")
	}

	return savePath, nil
}

var ErrLimitExceeded = errors.New("over limit")
var ErrAlreadyScrapedForFilter = errors.New("already scraped for filter")

func (i *UploadVideoFeature) SaveNewlyScraped(ctx context.Context, tx db.DBTX, jobId string, videoID string, filterID string, name string) error {
	queries := i.Querier.WithTx(tx)

	videoFromDb, err := queries.GetYtVideoWithJob(ctx, videoID)
	if err == nil && videoFromDb.FilterID.String == filterID {
		return ErrAlreadyScrapedForFilter
	} else if err == nil {
		// TODO connect to filter or job if multiple filters can exist
	} else if err != nil && err != sql.ErrNoRows {
		return err
	}

	job, err := queries.GetJob(ctx, jobId)
	if err != nil {
		return err
	}

	if job.VideosFound >= job.Limit.Int64 {
		return ErrLimitExceeded
	}

	_, err = queries.AddYtVideo(ctx, db.AddYtVideoParams{
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

var ErrHasDownloaded = errors.New("already downloaded")

func (i *UploadVideoFeature) SaveDownloadAttempt(ctx context.Context, tx db.DBTX, videoID string, savePath string, downloadError error) error {
	queries := i.Querier.WithTx(tx)
	attempts, err := queries.GetVideoWithDownloadAttempts(ctx, videoID)
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
		err := queries.AddDownloadAttempt(ctx, db.AddDownloadAttemptParams{
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
	err = queries.AddBlob(ctx, db.AddBlobParams{
		ID:   blobId.String(),
		Path: savePath,
	})
	if err != nil {
		return err
	}

	err = queries.AddDownloadAttempt(ctx, db.AddDownloadAttemptParams{
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
