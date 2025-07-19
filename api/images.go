package api

import (
	"context"
	"database/sql"

	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/api/util"
	"github.com/vkhobor/go-opencv/db"
)

type Picture struct {
	ID        string `json:"id"`
	BlobId    string `json:"blob_id"`
	YoutubeId string `json:"youtube_id"`
}

type Response struct {
	Pictures []Picture `json:"pictures"`
	Total    int       `json:"total"`
}

type ImagesRequest struct {
	Offset    int    `query:"offset"`
	Limit     int    `query:"limit"`
	YoutubeID string `query:"youtube_id"`
}

type HandleImagesResponse struct {
	Body Response
}

func HandleImages(sqlDB *sql.DB) util.Handler[ImagesRequest, HandleImagesResponse] {

	return func(ctx context.Context, e *ImagesRequest) (*HandleImagesResponse, error) {
		queries := db.New(sqlDB)
		res, err := queries.GetPictures(ctx, db.GetPicturesParams{
			Limit:               int64(e.Limit),
			Offset:              int64(e.Offset),
			IsFilterByYoutubeID: e.YoutubeID != "",
			YoutubeID: sql.NullString{
				String: e.YoutubeID,
				Valid:  true,
			},
		})
		if err != nil {
			return nil, err
		}

		count, err := queries.AllPicturesCount(ctx, db.AllPicturesCountParams{
			IsFilterByYoutubeID: e.YoutubeID != "",
			YoutubeID: sql.NullString{
				String: e.YoutubeID,
				Valid:  true,
			},
		})
		if err != nil {
			return nil, err
		}

		resp := Response{
			Total: int(count),
			Pictures: lo.Map(res, func(row db.GetPicturesRow, index int) Picture {
				return Picture{
					ID:        row.ID,
					BlobId:    row.BlobStorageID.String,
					YoutubeId: row.YtVideoID.String,
				}
			}),
		}

		return &HandleImagesResponse{Body: resp}, nil
	}
}
