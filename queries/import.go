package queries

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/db"
)

func (jc *Queries) GetRefImages(video DownlodedVideo) ([]string, error) {
	res, err := jc.Queries.GetFilterForJob(context.Background(), video.JobID)
	if err != nil {
		return []string{}, err
	}

	return lo.Map(res, func(item db.GetFilterForJobRow, i int) string {
		return item.Path
	}), nil
}

func (jc *Queries) StartImportAttempt(video DownlodedVideo) (string, error) {
	imported, err := jc.CheckImportedAlready(context.Background(), video.ID)
	if err != nil {
		return "", err
	}

	if imported {
		return "", ErrHasImported
	}

	importAttemptId := uuid.New().String()
	_, err = jc.Queries.AddImportAttempt(context.Background(), db.AddImportAttemptParams{
		ID: importAttemptId,
		YtVideoID: sql.NullString{
			String: video.ID,
			Valid:  true,
		},
		FilterID: sql.NullString{
			String: video.FilterID,
			Valid:  true,
		},
		Progress: sql.NullInt64{
			Int64: 0,
			Valid: true,
		},
		Error: sql.NullString{
			String: "",
			Valid:  false,
		},
	})

	if err != nil {
		return "", err
	}

	return importAttemptId, nil
}

func (jc *Queries) UpdateError(id string, err error) error {
	return jc.Queries.UpdateImportAttemptError(context.Background(), db.UpdateImportAttemptErrorParams{
		ID: id,
		Error: sql.NullString{
			String: err.Error(),
			Valid:  true,
		},
	})
}

func (jc *Queries) UpdateProgress(id string, progress int) error {
	if progress >= 100 {
		return ErrCannotUpdateTo100
	}

	return jc.updateProgress(id, progress)
}

func (jc *Queries) updateProgress(id string, progress int) error {
	return jc.Queries.UpdateImportAttemptProgress(context.Background(), db.UpdateImportAttemptProgressParams{
		ID: id,
		Progress: sql.NullInt64{
			Int64: int64(progress),
			Valid: true,
		},
	})
}

var ErrCannotUpdateTo100 = errors.New("cannot update to 100")
var ErrHasImported = errors.New("already imported")

func (jc *Queries) CheckImportedAlready(ctx context.Context, videoID string) (bool, error) {
	videos, err := jc.Queries.GetVideoWithImportAttempts(ctx, videoID)

	if err != nil {
		return false, err
	}

	if len(videos) > 0 {
		hasSuccessful := lo.SomeBy(videos, func(item db.GetVideoWithImportAttemptsRow) bool {
			return !item.Error.Valid && item.Progress.Int64 == 100
		})
		if hasSuccessful {
			return true, nil
		}
	}

	return false, nil
}

func (jc *Queries) FinishImport(video ImportedVideo, importAttemptId string) error {
	imported, err := jc.CheckImportedAlready(context.Background(), video.ID)
	if err != nil {
		return err
	}

	if imported {
		return ErrHasImported
	}

	err = jc.updateProgress(importAttemptId, 100)

	if err != nil {
		return err
	}

	for _, frame := range video.ExtractedFrames {
		blobID := uuid.New()
		_ = jc.Queries.AddBlob(context.Background(), db.AddBlobParams{
			ID:   blobID.String(),
			Path: frame.Path,
		})
		_ = jc.Queries.AddPicture(context.Background(), db.AddPictureParams{
			ID: uuid.New().String(),
			ImportAttemptID: sql.NullString{
				String: importAttemptId,
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
	return nil
}
