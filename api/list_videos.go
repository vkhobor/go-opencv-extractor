package api

import (
	"context"

	u "github.com/vkhobor/go-opencv/api/util"
	"github.com/vkhobor/go-opencv/db"
	"github.com/vkhobor/go-opencv/queries"
)

type ListVideoBody struct {
	VideoID  string `json:"video_id"`
	Name     string `json:"name"`
	Progress int    `json:"progress"`
}

type ListVideosResponse struct {
	Body []ListVideoBody
}

func HandleListVideos(dbQ *db.Queries) u.Handler[struct{}, ListVideosResponse] {
	return func(ctx context.Context, e *struct{}) (*ListVideosResponse, error) {
		queries := queries.Queries{
			Queries: dbQ,
		}
		res := queries.GetDownloadedVideos(true)

		videosResponse := []ListVideoBody{}
		for _, video := range res {
			videosResponse = append(videosResponse, ListVideoBody{
				VideoID:  video.ID,
				Progress: int(video.ImportProgress),
				Name:     video.Name,
			})
		}

		return &ListVideosResponse{
			Body: videosResponse,
		}, nil
	}
}
