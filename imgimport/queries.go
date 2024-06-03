package imgimport

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/db"
	"github.com/vkhobor/go-opencv/download"
)

type Queries struct {
	Queries *db.Queries
}

func (jc *Queries) GetRefImages(video download.DownlodedVideo) ([]string, error) {
	res, err := jc.Queries.GetFilterForJob(context.Background(), video.JobID)
	if err != nil {
		return []string{}, err
	}

	return lo.Map(res, func(item db.GetFilterForJobRow, i int) string {
		return item.Path
	}), nil
}

func (jc *Queries) StartImportAttempt(video download.DownlodedVideo) (string, error) {
	importAttemptId := uuid.New().String()
	_, err := jc.Queries.AddImportAttempt(context.Background(), db.AddImportAttemptParams{
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

func (jc *Queries) UpdateProgress(id string, progress int) {
	// TODO: Implement
	return
}

func (jc *Queries) SaveFrames(video ImportedVideo, importAttemptId string) {
	for _, frame := range video.ExtractedFrames {
		blobID := uuid.New()
		_, _ = jc.Queries.AddBlob(context.Background(), db.AddBlobParams{
			ID:   blobID.String(),
			Path: frame.Path,
		})
		_, _ = jc.Queries.AddPicture(context.Background(), db.AddPictureParams{
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
}
