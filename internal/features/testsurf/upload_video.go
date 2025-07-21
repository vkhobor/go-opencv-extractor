package testsurf

import (
	"context"
	"io"
	"os"

	"github.com/vkhobor/go-opencv/internal/config"
	"github.com/vkhobor/go-opencv/internal/mlog"
	"github.com/vkhobor/go-opencv/internal/video"
)

type UploadVideoFeature struct {
	Config config.DirectoryConfig
}

func (f *UploadVideoFeature) UploadVideo(ctx context.Context, videoData io.Reader) error {
	mlog.Log().Info("Uploading test video")

	filePath := f.Config.GetTestVideoPath()

	if _, err := os.Stat(filePath); err == nil {
		mlog.Log().Info("Removing existing test video file", "path", filePath)
		if err := os.Remove(filePath); err != nil {
			mlog.Log().Error("Failed to remove existing video file", "error", err, "path", filePath)
			return err
		}
	}

	dst, err := os.Create(filePath)
	if err != nil {
		mlog.Log().Error("Failed to create video file", "error", err, "path", filePath)
		return err
	}
	defer dst.Close()

	mlog.Log().Info("Copying video data to file", "path", filePath)

	written, err := io.Copy(dst, videoData)

	if err != nil {
		mlog.Log().Error("Failed to copy video data", "error", err, "path", filePath)
		return err
	}

	mlog.Log().Info("Successfully saved test video", "path", filePath, "bytesWritten", written)

	cachedTestVideoExtractor.Close()
	cachedTestVideoExtractor, err = video.NewFrameExtractor(filePath)
	if err != nil {
		return err
	}

	return nil
}
