package filters

import (
	"context"

	u "github.com/vkhobor/go-opencv/api/util"
	"github.com/vkhobor/go-opencv/db"
	"github.com/vkhobor/go-opencv/features"
)

type ReferenceGetRequest struct {
	ID string `path:"id"`
}

type ReferenceGetResponse struct {
	Body features.ReferenceGetFeatureResponse `json:"body"`
}

func HandleReferenceGet(queries *db.Queries) u.Handler[ReferenceGetRequest, ReferenceGetResponse] {
	feature := &features.ReferenceGetFeature{
		Queries: queries,
	}

	return func(ctx context.Context, req *ReferenceGetRequest) (*ReferenceGetResponse, error) {
		res, err := feature.GetReference(ctx)
		if err != nil {
			return nil, err
		}

		return &ReferenceGetResponse{Body: res}, nil
	}
}
