package import_video

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/internal/db"
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

func (d *ImportVideoFeature) getRefImages(ctx context.Context, tx db.DBTX, jobId string) (FilterWithPaths, error) {
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

func (jc *ImportVideoFeature) startImportAttempt(ctx context.Context, tx db.DBTX, videoID string, filterID string) (string, error) {
	queries := jc.Querier.WithTx(tx)

	importAttemptId := uuid.New().String()
	_, err := queries.AddImportAttempt(ctx, db.AddImportAttemptParams{
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

func (jc *ImportVideoFeature) updateError(ctx context.Context, tx db.DBTX, id string, err error) error {
	queries := jc.Querier.WithTx(tx)
	return queries.UpdateImportAttemptError(ctx, db.UpdateImportAttemptErrorParams{
		ID: id,
		Error: sql.NullString{
			String: err.Error(),
			Valid:  true,
		},
	})
}

func (jc *ImportVideoFeature) updateProgress(ctx context.Context, id string, progress int) error {
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

func (jc *ImportVideoFeature) addFrameToVideo(ctx context.Context, frame Frame, importAttemptId string) error {
	tx, err := jc.SqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queries := jc.Querier.WithTx(tx)

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

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (jc *ImportVideoFeature) finishImport(ctx context.Context, tx db.DBTX, videoID string, importAttemptId string) error {
	queries := jc.Querier.WithTx(tx)

	err := queries.UpdateImportAttemptProgress(ctx, db.UpdateImportAttemptProgressParams{
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
