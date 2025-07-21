package features

import (
	"context"
	"database/sql"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/vkhobor/go-opencv/internal/config"
	"github.com/vkhobor/go-opencv/internal/db"
)

type ReferenceUploadFeature struct {
	SqlDB   TXer
	Querier QuerierWithTx
	Config  config.DirectoryConfig
}

type ReferenceConfig struct {
	RatioTestThreshold         float64
	MinThresholdForSURFMatches float64
	MinSURFMatches             int64
	MseSkip                    float64
}

// defaultFilterId is temporary until dynamic filters are implemented
const (
	defaultFilterId   = "1fed33d4-0ea3-4b84-909c-261e4b2a3d43"
	defaultFilterType = "SURF"
)

func (f *ReferenceUploadFeature) UploadReference(
	ctx context.Context,
	file io.Reader,
	fileName string,
	config ReferenceConfig) error {
	path, err := f.overrideReferencesOnDisk(file, fileName)
	if err != nil {
		return err
	}

	tx, err := f.SqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = f.saveToDb(ctx, tx, path, config)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (f *ReferenceUploadFeature) overrideReferencesOnDisk(file io.Reader, fileName string) (string, error) {
	err := os.MkdirAll(f.Config.GetReferencesDir(), os.ModePerm)
	if err != nil {
		return "", err
	}

	err = os.RemoveAll(f.Config.GetReferencesDir())
	if err != nil {
		return "", err
	}

	path := filepath.Join(f.Config.GetReferencesDir(), fileName)
	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return path, nil
}

func (f *ReferenceUploadFeature) saveToDb(ctx context.Context, tx db.DBTX, path string, config ReferenceConfig) error {
	queries := f.Querier.WithTx(tx)

	err := f.upsertFilter(ctx, tx, config)
	if err != nil {
		return err
	}

	blobId := uuid.NewString()
	err = queries.AddBlob(ctx, db.AddBlobParams{
		ID:   blobId,
		Path: path,
	})
	if err != nil {
		return err
	}

	err = queries.DeleteImagesOnFilter(ctx, sql.NullString{
		String: defaultFilterId,
		Valid:  true,
	})
	if err != nil {
		return err
	}

	_, err = queries.AttachImageToFilter(ctx, db.AttachImageToFilterParams{
		FilterID: sql.NullString{
			String: defaultFilterId,
			Valid:  true,
		},
		BlobStorageID: sql.NullString{
			String: blobId,
			Valid:  true,
		},
	})
	return err
}

func (f *ReferenceUploadFeature) upsertFilter(ctx context.Context, tx db.DBTX, config ReferenceConfig) error {
	queries := f.Querier.WithTx(tx)

	_, err := queries.AddFilter(ctx, db.AddFilterParams{
		ID: defaultFilterId,
		Name: sql.NullString{
			String: "Default",
			Valid:  true,
		},
		Discriminator: sql.NullString{
			String: defaultFilterType,
			Valid:  true,
		},
		Ratiotestthreshold: sql.NullFloat64{
			Float64: config.RatioTestThreshold,
			Valid:   true,
		},
		Minthresholdforsurfmatches: sql.NullFloat64{
			Float64: config.MinThresholdForSURFMatches,
			Valid:   true,
		},
		Minsurfmatches: sql.NullInt64{
			Int64: config.MinSURFMatches,
			Valid: true,
		},
		Mseskip: sql.NullFloat64{
			Float64: config.MseSkip,
			Valid:   true,
		},
	})
	return err
}
