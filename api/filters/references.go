package filters

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/samber/lo"
	u "github.com/vkhobor/go-opencv/api/util"
	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/db"
)

const (
	megabyte = 1 << 20 // 1 megabyte = 2^20 bytes
)

type ReferenceUploadRequest struct {
	RawBody huma.MultipartFormFiles[struct {
		File                       huma.FormFile `form:"file" required:"true"`
		RatioTestThreshold         float64       `form:"ratioTestThreshold" required:"true"`
		MinThresholdForSURFMatches float64       `form:"minThresholdForSURFMatches" required:"true"`
		MinSURFMatches             int64         `form:"minSURFMatches" required:"true"`
		MseSkip                    float64       `form:"mseSkip" required:"true"`
	}]
}

type ReferenceUploadResponse struct {
	Status string `json:"status"`
}

// TODO migrate to fully dynamic filters
var defaultFilterId = "1fed33d4-0ea3-4b84-909c-261e4b2a3d43"

func HandleReferenceUpload(queries *db.Queries, config config.DirectoryConfig) u.Handler[ReferenceUploadRequest, ReferenceUploadResponse] {
	return func(ctx context.Context, req *ReferenceUploadRequest) (*ReferenceUploadResponse, error) {
		data := req.RawBody.Data()

		path, err := saveToDisk(data.File.File, config, data.File.Filename)
		if err != nil {
			return nil, err
		}

		err = saveToDb(ctx, queries, path,
			data.RatioTestThreshold,
			data.MinThresholdForSURFMatches,
			data.MinSURFMatches,
			data.MseSkip)
		if err != nil {
			return nil, err
		}

		return &ReferenceUploadResponse{Status: "created"}, nil
	}
}

func saveToDb(
	ctx context.Context,
	queries *db.Queries,
	path string,
	ratioTestThreshold float64,
	minThresholdForSURFMatches float64,
	minSURFMatches int64,
	mseSkip float64) error {
	id := defaultFilterId

	filters, err := queries.GetFilters(ctx)
	if err != nil {
		return err
	}

	exists := lo.SomeBy(filters, func(item db.GetFiltersRow) bool {
		return item.ID == id
	})
	slog.Debug("Filter exists", "exists", exists, "filters", filters, "id", id)

	if !exists {
		_, err := queries.AddFilter(ctx, db.AddFilterParams{
			ID: id,
			Name: sql.NullString{
				String: "Default",
				Valid:  true,
			},
			Discriminator: sql.NullString{
				String: "SURF",
				Valid:  true,
			},
			Ratiotestthreshold: sql.NullFloat64{
				Float64: ratioTestThreshold,
				Valid:   true,
			},
			Minthresholdforsurfmatches: sql.NullFloat64{
				Float64: minThresholdForSURFMatches,
				Valid:   true,
			},
			Minsurfmatches: sql.NullInt64{
				Int64: minSURFMatches,
				Valid: true,
			},
			Mseskip: sql.NullFloat64{
				Float64: mseSkip,
				Valid:   true,
			},
		})

		if err != nil {
			return err
		}
	}

	blobId := uuid.NewString()
	err = queries.AddBlob(ctx, db.AddBlobParams{
		ID:   blobId,
		Path: path,
	})
	if err != nil {
		return err
	}

	_, err = queries.AttachImageToFilter(ctx, db.AttachImageToFilterParams{
		FilterID: sql.NullString{
			String: id,
			Valid:  true,
		},
		BlobStorageID: sql.NullString{
			String: blobId,
			Valid:  true,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func saveToDisk(file io.Reader, config config.DirectoryConfig, fileName string) (string, error) {
	err := os.MkdirAll(config.GetReferencesDir(), os.ModePerm)
	if err != nil {
		return "", err
	}

	dst, err := os.Create(fmt.Sprintf("%s/%s", config.GetReferencesDir(), fileName))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	path := fmt.Sprintf("%s/%s", config.GetReferencesDir(), fileName)
	return path, nil
}

func HandleGetReferences(queries *db.Queries) http.HandlerFunc {
	type reference struct {
		ID string `json:"id"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			res, err := queries.GetFilters(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			response := lo.FilterMap(res, func(item db.GetFiltersRow, index int) (reference, bool) {
				return reference{
					ID: item.BlobID.String,
				}, item.BlobID.Valid
			})

			render.JSON(w, r, response)
		},
	)
}

func HandleDeleteAllReferences(queries *db.Queries) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := queries.DeleteImagesOnFilter(r.Context(), sql.NullString{
				String: defaultFilterId,
				Valid:  true,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			render.Status(r, http.StatusNoContent)
		},
	)
}
