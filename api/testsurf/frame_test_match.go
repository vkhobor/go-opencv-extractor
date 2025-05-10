package testsurf

import (
	"context"

	u "github.com/vkhobor/go-opencv/api/util"
)

type FrameMatchingTestRequest struct {
	FrameNum           int     `query:"framenum" required:"true"`
	RatioCheck         float64 `query:"ratiocheck" required:"true"`
	MinMatches         int     `query:"minmatches" required:"true"`
	GoodMatchThreshold float64 `query:"goodmatchthreshold" required:"true"`
}

type FrameMatchingTestBody struct {
	Matched bool `json:"matched"`
}

type FrameMatchingTestResponse struct {
	Body FrameMatchingTestBody `json:"body"`
}

func HandleFrameMatchingTest() u.Handler[FrameMatchingTestRequest, FrameMatchingTestResponse] {
	return func(ctx context.Context, req *FrameMatchingTestRequest) (*FrameMatchingTestResponse, error) {
		return &FrameMatchingTestResponse{
			Body: FrameMatchingTestBody{
				Matched: false,
			},
		}, nil
	}
}
