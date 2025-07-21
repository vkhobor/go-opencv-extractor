package features

import (
	"context"

	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/internal/db"
)

type ReferenceGetFeature struct {
	Querier QuerierWithTx
}

type ReferenceGetFeatureResponse struct {
	ID                         string
	Name                       string
	Discriminator              string
	Ratiotestthreshold         float64
	Minthresholdforsurfmatches float64
	Minsurfmatches             int64
	Mseskip                    float64
	BlobIds                    []string
}

func (f *ReferenceGetFeature) GetReference(ctx context.Context) (ReferenceGetFeatureResponse, error) {
	res, err := f.Querier.GetFilterById(ctx, defaultFilterId)
	if err != nil {
		return ReferenceGetFeatureResponse{}, err
	}

	if len(res) == 0 {
		return ReferenceGetFeatureResponse{}, nil
	}

	return ReferenceGetFeatureResponse{
		BlobIds: lo.Map(res, func(path db.GetFilterByIdRow, index int) string {
			return path.BlobStorageID.String
		}),
		ID:                         res[0].ID,
		Name:                       res[0].Name.String,
		Discriminator:              res[0].Discriminator.String,
		Ratiotestthreshold:         res[0].Ratiotestthreshold.Float64,
		Minthresholdforsurfmatches: res[0].Minthresholdforsurfmatches.Float64,
		Minsurfmatches:             res[0].Minsurfmatches.Int64,
		Mseskip:                    res[0].Mseskip.Float64,
	}, nil
}
