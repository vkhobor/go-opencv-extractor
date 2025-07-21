package api

import (
	"context"
	"database/sql"
	"log/slog"

	u "github.com/vkhobor/go-opencv/internal/api/util"
	"github.com/vkhobor/go-opencv/internal/db"
)

type ListVideoBody struct {
	VideoID  string `json:"video_id"`
	Name     string `json:"name"`
	Progress int    `json:"progress"`
}

type ListVideosResponse struct {
	Body []ListVideoBody
}

func HandleListVideos(dbQ *sql.DB) u.Handler[struct{}, ListVideosResponse] {
	return func(ctx context.Context, e *struct{}) (*ListVideosResponse, error) {
		queries := db.New(dbQ)
		val, err := queries.GetVideosDownloaded(ctx)
		if err != nil {
			slog.Error("GetDownloadedVideos: Error while getting downloaded videos", "error", err)
			return nil, err
		}

		videosResponse := []ListVideoBody{}
		for _, video := range val {
			videosResponse = append(videosResponse, ListVideoBody{
				VideoID:  video.YtVideoID,
				Progress: int(video.ImportProgress),
				Name:     video.YtVideoName.String,
			})
		}

		return &ListVideosResponse{
			Body: videosResponse,
		}, nil
	}
}
