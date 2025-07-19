package import_video

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/db"
)

type Frame struct {
	FrameNumber int
	Path        string
}

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

func (d *ImportVideoFeature) GetRefImages(ctx context.Context, tx db.DBTX, jobId string) (FilterWithPaths, error) {
	queries := d.Querier.WithTx(tx)

	res, err := queries.GetFilterForJob(ctx, jobId)
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

func (jc *ImportVideoFeature) StartImportAttempt(ctx context.Context, tx db.DBTX, videoID string, filterID string) (string, error) {
	queries := jc.Querier.WithTx(tx)

	imported, err := jc.CheckImportedAlready(ctx, tx, videoID)
	if err != nil {
		return "", err
	}

	if imported {
		return "", ErrHasImported
	}

	importAttemptId := uuid.New().String()
	_, err = queries.AddImportAttempt(ctx, db.AddImportAttemptParams{
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

func (jc *ImportVideoFeature) UpdateError(ctx context.Context, tx db.DBTX, id string, err error) error {
	queries := jc.Querier.WithTx(tx)
	return queries.UpdateImportAttemptError(ctx, db.UpdateImportAttemptErrorParams{
		ID: id,
		Error: sql.NullString{
			String: err.Error(),
			Valid:  true,
		},
	})
}

func (jc *ImportVideoFeature) UpdateProgress(ctx context.Context, tx db.DBTX, id string, progress int) error {
	if progress >= 100 {
		return ErrCannotUpdateTo100
	}

	return jc.Querier.UpdateImportAttemptProgress(ctx, db.UpdateImportAttemptProgressParams{
		ID: id,
		Progress: sql.NullInt64{
			Int64: int64(progress),
			Valid: true,
		},
	})
}

var ErrCannotUpdateTo100 = errors.New("cannot update to 100")
var ErrHasImported = errors.New("already imported")

func (jc *ImportVideoFeature) CheckImportedAlready(ctx context.Context, tx db.DBTX, videoID string) (bool, error) {
	queries := jc.Querier.WithTx(tx)
	videos, err := queries.GetVideoWithImportAttempts(ctx, videoID)

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

func (jc *ImportVideoFeature) AddFrameToVideo(ctx context.Context, tx db.DBTX, videoID string, frame Frame, importAttemptId string) error {
	queries := jc.Querier.WithTx(tx)
	if imported, err := jc.CheckImportedAlready(ctx, tx, videoID); err != nil {
		return err
	} else if imported {
		return ErrHasImported
	}

	blobID := uuid.New()
	_ = queries.AddBlob(ctx, db.AddBlobParams{
		ID:   blobID.String(),
		Path: frame.Path,
	})
	_ = queries.AddPicture(ctx, db.AddPictureParams{
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

func (jc *ImportVideoFeature) FinishImport(ctx context.Context, tx db.DBTX, videoID string, importAttemptId string) error {
	queries := jc.Querier.WithTx(tx)

	imported, err := jc.CheckImportedAlready(ctx, tx, videoID)
	if err != nil {
		return err
	}

	if imported {
		return ErrHasImported
	}

	err = queries.UpdateImportAttemptProgress(ctx, db.UpdateImportAttemptProgressParams{
		ID: importAttemptId,
		Progress: sql.NullInt64{
			Int64: int64(100),
			Valid: true,
		},
	})

	if err != nil {
		return err
	}

	return nil
}
