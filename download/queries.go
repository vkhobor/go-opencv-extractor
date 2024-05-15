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

func (jc *Queries) DownloadSaved(video DownlodedVideo) {
	slog.Debug("Saving downloaded", "video", video)
	if video.Error != nil {
		_, err := jc.Queries.UpdateStatus(context.Background(), db.UpdateStatusParams{
			ID: video.ID,
			Status: sql.NullString{
				String: "errored",
				Valid:  true,
			},
			Error: sql.NullString{
				Valid:  true,
				String: video.Error.Error(),
			},
		})
		if err != nil {
			slog.Error("Error while updating status", "error", err)
		}

		return
	}
	blobId := uuid.New()
	_, errAddBlob := jc.Queries.AddBlob(context.Background(), db.AddBlobParams{
		ID:   blobId.String(),
		Path: video.SavePath,
	})

	_, errUpdateStatus := jc.Queries.UpdateStatus(context.Background(), db.UpdateStatusParams{
		ID: video.ID,
		Status: sql.NullString{
			String: "downloaded",
			Valid:  true,
		},
	})

	jc.Queries.AddBlobToVideo(context.Background(), db.AddBlobToVideoParams{
		BlobStorageID: sql.NullString{
			String: blobId.String(),
			Valid:  true,
		},
		ID: video.ID,
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
