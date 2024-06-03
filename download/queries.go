package download

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/vkhobor/go-opencv/db"
	"github.com/vkhobor/go-opencv/scraper"
)

type Queries struct {
	Queries *db.Queries
}

func (jc *Queries) GetDownloadedVideos() []DownlodedVideo {
	val, err := jc.Queries.GetVideosDownloaded(context.Background())
	if err != nil {
		return []DownlodedVideo{}
	}
	result := make([]DownlodedVideo, len(val))
	for i, v := range val {
		result[i] = DownlodedVideo{ScrapedVideo: scraper.ScrapedVideo{ID: v.ID}, SavePath: v.Path}
	}
	return result
}

func (jc *Queries) SaveDownloadAttempt(video DownlodedVideo) {
	slog.Debug("Saving downloaded", "video", video)
	if video.Error != nil {

	}
	blobId := uuid.New()
	_, errAddBlob := jc.Queries.AddBlob(context.Background(), db.AddBlobParams{
		ID:   blobId.String(),
		Path: video.SavePath,
	})

	_, errUpdateStatus := jc.Queries.AddDownloadAttempt(context.Background(), db.AddDownloadAttemptParams{
		YtVideoID: sql.NullString{
			String: video.ID,
			Valid:  true,
		},
		Error: sql.NullString{
			String: video.Error.Error(),
			Valid:  video.Error != nil,
		},
		BlobStorageID: sql.NullString{
			String: blobId.String(),
			Valid:  true,
		},
	})

	if errAddBlob != nil || errUpdateStatus != nil {
		slog.Error("Error while updating status", "error", errAddBlob, "error2", errUpdateStatus)
		RemoveAllPaths(video.SavePath)
		return
	}
}

func RemoveAllPaths(files ...string) {
	for _, file := range files {
		_ = os.Remove(file)
	}
}
