package testsurf

import (
	"context"

	"github.com/vkhobor/go-opencv/mlog"
	"github.com/vkhobor/go-opencv/video/videoiter"
	"gocv.io/x/gocv"
)

// VideoMetadataFeature handles retrieving metadata for test videos
type VideoMetadataFeature struct {
	// Add any dependencies here
}

// VideoMetadata contains basic information about a video
type VideoMetadata struct {
	MaxFrame int
	FPS      float64
	Width    int
	Height   int
	Duration float64 // in seconds
}

// GetVideoMetadata retrieves metadata for the currently uploaded test video
func (f *VideoMetadataFeature) GetVideoMetadata(ctx context.Context) (VideoMetadata, error) {
	mlog.Log().Info("Retrieving video metadata")

	// TODO: Implement actual video metadata retrieval
	// This would typically involve:
	// 1. Finding the current test video
	// 2. Opening it with gocv
	// 3. Extracting metadata
	
	// For now, this is a dummy implementation
	// In a real implementation, you would:
	// video, err := videoiter.NewVideo(videoPath)
	// if err != nil {
	//     return VideoMetadata{}, err
	// }
	// defer video.Close()
	
	// And then extract properties like:
	// metadata.FPS = video.FPS()
	// metadata.MaxFrame = int(video.FrameCount())
	// etc.
	
	// Mock implementation returning dummy data
	return VideoMetadata{
		MaxFrame: 0, // Zero frames indicates no video loaded
		FPS:      0,
		Width:    0,
		Height:   0,
		Duration: 0,
	}, nil
}

// GetFrameCount gets the number of frames in a video file
func (f *VideoMetadataFeature) GetFrameCount(videoPath string) (int, error) {
	video, err := gocv.OpenVideoCapture(videoPath)
	if err != nil {
		mlog.Log().Error("Failed to open video", "error", err)
		return 0, err
	}
	defer video.Close()

	// Get video properties
	frameCount := int(video.Get(gocv.VideoCaptureFrameCount))
	mlog.Log().Info("Video frame count", "count", frameCount)
	
	return frameCount, nil
}

// GetVideoFromPath opens a video for inspection
func (f *VideoMetadataFeature) GetVideoFromPath(videoPath string) (*videoiter.Video, error) {
	return videoiter.NewVideo(videoPath)
}