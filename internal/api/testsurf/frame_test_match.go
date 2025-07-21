package testsurf

import (
	"context"

	u "github.com/vkhobor/go-opencv/internal/api/util"
	"github.com/vkhobor/go-opencv/internal/features/testsurf"
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
		feat := testsurf.FrameMatchingTestFeature{}

		ok, err := feat.TestFrameMatch(ctx, req.FrameNum, req.RatioCheck, req.MinMatches, req.GoodMatchThreshold)
		if err != nil {
			return nil, err
		}

		return &FrameMatchingTestResponse{
			Body: FrameMatchingTestBody{
				Matched: ok,
			},
		}, nil
	}
}
