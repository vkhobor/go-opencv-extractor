package queries

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/db"
)

type FilterWithPaths struct {
	ID                         string
	Name                       string
	Discriminator              string
	RatioTestThreshold         float64
	MinThresholdForSURFMatches float64
	MinSURFMatches             int64
	MSESkip                    float64
	Paths                      []string
}

func (jc *Queries) GetRefImages(jobId string) (FilterWithPaths, error) {
	res, err := jc.Queries.GetFilterForJob(context.Background(), jobId)
	if err != nil {
		return FilterWithPaths{}, err
	}

	if len(res) == 0 {
		return FilterWithPaths{}, nil
	}

	first := res[0]
	return FilterWithPaths{
		ID:                         first.ID,
		Name:                       first.Name.String,
		Discriminator:              first.Discriminator.String,
		RatioTestThreshold:         first.Ratiotestthreshold.Float64,
		MinThresholdForSURFMatches: first.Minthresholdforsurfmatches.Float64,
		MinSURFMatches:             first.Minsurfmatches.Int64,
		MSESkip:                    first.Mseskip.Float64,
		Paths: lo.Map(res, func(item db.GetFilterForJobRow, _ int) string {
			return item.Path
		}),
	}, nil
}

func (jc *Queries) StartImportAttempt(videoID string, filterID string) (string, error) {
	imported, err := jc.CheckImportedAlready(context.Background(), videoID)
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
			String: videoID,
			Valid:  true,
		},
		FilterID: sql.NullString{
			String: filterID,
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

func (jc *Queries) AddFrameToVideo(videoID string, frame Frame, importAttemptId string) error {
	if imported, err := jc.CheckImportedAlready(context.Background(), videoID); err != nil {
		return err
	} else if imported {
		return ErrHasImported
	}

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

	return nil
}

func (jc *Queries) FinishImport(videoID string, importAttemptId string) error {
	imported, err := jc.CheckImportedAlready(context.Background(), videoID)
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

	return nil
}
