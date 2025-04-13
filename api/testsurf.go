package api

import (
	"context"

	"github.com/vkhobor/go-opencv/api/util"
	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/db"
)

type UploadVideoRequest struct {
	Body struct {
		VideoData string `json:"video_data"`
	} `json:"body"`
}

type UploadVideoResponse struct {
	Message string `json:"message"`
}

func HandleTestSurfUpload(queries *db.Queries, config config.ServerConfig) util.Handler[UploadVideoRequest, UploadVideoResponse] {
	return func(ctx context.Context, req *UploadVideoRequest) (*UploadVideoResponse, error) {
		// Logic to handle video upload and delete the previous video
		return &UploadVideoResponse{Message: "Video uploaded and previous video deleted"}, nil
	}
}

type TestSurfCheckRequest struct {
	FrameNum           int     `query:"framenum"`
	RatioCheck         float64 `query:"ratiocheck"`
	MinMatches         int     `query:"minmatches"`
	GoodMatchThreshold float64 `query:"goodmatchthreshold"`
}

type TestSurfCheckResponse struct {
	Matched bool `json:"matched"`
}

func HandleTestSurfCheck(queries *db.Queries) util.Handler[TestSurfCheckRequest, TestSurfCheckResponse] {
	return func(ctx context.Context, req *TestSurfCheckRequest) (*TestSurfCheckResponse, error) {
		// Logic to check matching based on query parameters
		return &TestSurfCheckResponse{Matched: true}, nil
	}
}

type FrameRequest struct {
	FrameNum int `query:"framenum"`
}

type FrameResponse struct {
	ImageData string `json:"image_data"`
}

func HandleTestSurfFrame(queries *db.Queries) util.Handler[FrameRequest, FrameResponse] {
	return func(ctx context.Context, req *FrameRequest) (*FrameResponse, error) {
		// Logic to serve an image for a specific frame number
		return &FrameResponse{ImageData: "Image data for frame"}, nil
	}
}

type MaxFrameResponse struct {
	MaxFrame int `json:"maxframe"`
}

func HandleTestSurfMaxFrame(queries *db.Queries) util.Handler[struct{}, MaxFrameResponse] {
	return func(ctx context.Context, req *struct{}) (*MaxFrameResponse, error) {
		// Logic to return max frame
		return &MaxFrameResponse{MaxFrame: 250000}, nil
	}
}
