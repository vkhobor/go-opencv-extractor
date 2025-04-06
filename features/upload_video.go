package features

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/db"
	"github.com/vkhobor/go-opencv/mlog"
	"github.com/vkhobor/go-opencv/queries"
)

type UploadVideoFeature struct {
	Queries  *queries.Queries
	Config   config.DirectoryConfig
	WakeJobs chan<- struct{}
}

func (i *UploadVideoFeature) DownloadVideo(data io.Reader, filterId string, name string) (savePath string, error error) {
	jobID := uuid.New().String()
	mlog.Log().Info("Creating new job", "jobID", jobID, "filterID", filterId)
	_, err := i.Queries.Queries.CreateJob(context.Background(), db.CreateJobParams{
		FilterID: sql.NullString{
			String: filterId,
			Valid:  true,
		},
		Limit: sql.NullInt64{
			Valid: true,
			Int64: 1,
		},
		ID: jobID,
	})

	videoId := uuid.New().String()
	mlog.Log().Info("Saving newly scraped video", "jobID", jobID, "videoID", videoId, "filterID", filterId)
	err = i.Queries.SaveNewlyScraped(jobID, videoId, filterId, name)
	if err != nil {
		mlog.Log().Error("Failed to save newly scraped video", "error", err)
		return "", err
	}

	folderPath := i.Config.GetVideosDir()
	mlog.Log().Debug("Creating videos directory", "path", folderPath)
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		mlog.Log().Error("Failed to create videos directory", "error", err)
		return "", err
	}

	id := uuid.New()
	fileName := fmt.Sprintf("%v_%v.mp4", id.String(), videoId)
	filePath := filepath.Join(folderPath, fileName)
	mlog.Log().Info("Preparing to save video file", "path", filePath)

	savePath = filePath
	dst, err := os.Create(filePath)
	if err != nil {
		mlog.Log().Error("Failed to create video file", "error", err, "path", filePath)
		return "", err
	}
	defer dst.Close()

	mlog.Log().Info("Starting video data copy", "videoID", videoId)
	buffer := make([]byte, 1024)
	bytesWritten := 0
	for {
		n, err := data.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			mlog.Log().Error("Error reading video data", "error", err)
			return "", err
		}
		dst.Write(buffer[:n])
		bytesWritten += n
	}
	mlog.Log().Info("Completed video data copy", "videoID", videoId, "bytesWritten", bytesWritten)

	err = i.Queries.SaveDownloadAttempt(videoId, filePath, err)
	if err != nil {
		mlog.Log().Error("Error while saving download attempt", "error", err, "videoID", videoId)
		return savePath, err
	}

	mlog.Log().Info("Successfully completed video download", "videoID", videoId, "path", filePath)

	select {
	case i.WakeJobs <- struct{}{}:
		mlog.Log().Info("Waking up jobs")
	default:
		mlog.Log().Info("Jobs already awake")
	}

	return savePath, nil
}
