package filters

import (
	"context"
	"database/sql"

	u "github.com/vkhobor/go-opencv/api/util"
	"github.com/vkhobor/go-opencv/features"
)

type ReferenceGetRequest struct {
	ID string `path:"id"`
}

type ReferenceGetResponse struct {
	Body features.ReferenceGetFeatureResponse `json:"body"`
}

func HandleReferenceGet(queries *sql.DB) u.Handler[ReferenceGetRequest, ReferenceGetResponse] {
	dbAdapter := u.NewDbAdapter(queries)

	feature := &features.ReferenceGetFeature{
		Querier: dbAdapter.Querier,
	}

	return func(ctx context.Context, req *ReferenceGetRequest) (*ReferenceGetResponse, error) {
		res, err := feature.GetReference(ctx)
		if err != nil {
			return nil, err
		}

		return &ReferenceGetResponse{Body: res}, nil
	}
}
