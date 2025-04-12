package filters

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/db"
)

const (
	megabyte = 1 << 20 // 1 megabyte = 2^20 bytes
)

// TODO migrate to fully dynamic filters
var defaultFilterId = "1fed33d4-0ea3-4b84-909c-261e4b2a3d43"

func HandleReferenceUpload(queries *db.Queries, config config.DirectoryConfig) http.HandlerFunc {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseMultipartForm(10 * megabyte)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ratioTestThreshold, err := strconv.ParseFloat(r.FormValue("ratioTestThreshold"), 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			minThresholdForSURFMatches, err := strconv.ParseFloat(r.FormValue("minThresholdForSURFMatches"), 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			minSURFMatches, err := strconv.ParseInt(r.FormValue("minSURFMatches"), 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			mseSkip, err := strconv.ParseFloat(r.FormValue("mseSkip"), 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			files := r.MultipartForm.File
			for _, headers := range files {
				for _, header := range headers {
					file, err := header.Open()
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					defer file.Close()

					path, err := saveToDisk(file, config, header.Filename)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					err = saveToDb(
						r.Context(),
						queries,
						path,
						ratioTestThreshold,
						minThresholdForSURFMatches,
						minSURFMatches,
						mseSkip)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
			}
			render.Status(r, http.StatusCreated)
		},
	)
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
