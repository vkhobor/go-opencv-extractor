package imgimport

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/db"
)

type Queries struct {
	Queries *db.Queries
}

func (jc *Queries) GetRefImages() ([]string, error) {
	val, err := jc.Queries.GetReferences(context.Background())
	if err != nil {
		return nil, err
	}

	return lo.Map(val, func(item db.GetReferencesRow, i int) string {
		return item.Path
	}), nil
}

func (jc *Queries) SaveImported(video ImportedVideo) {
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
