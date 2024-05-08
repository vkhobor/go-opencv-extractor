package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/db_sql"
)

func HandleImages(queries *db_sql.Queries) http.HandlerFunc {
	type picture struct {
		ID        string `json:"id"`
		BlobId    string `json:"blob_id"`
		YoutubeId string `json:"youtube_id"`
	}

	type response struct {
		Pictures []picture `json:"pictures"`
		Total    int       `json:"total"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			offset := r.URL.Query().Get("offset")
			limit := r.URL.Query().Get("limit")

			if offset == "" || limit == "" {
				render.Status(r, http.StatusBadRequest)
				render.PlainText(w, r, "No offset or limit provided")
				return
			}

			limitAsInt, err := strconv.Atoi(limit)
			if err != nil {
				render.Status(r, http.StatusBadRequest)
				render.PlainText(w, r, "Invalid limit")
				return
			}

			offsetAsInt, err := strconv.Atoi(offset)
			if err != nil {
				render.Status(r, http.StatusBadRequest)
				render.PlainText(w, r, "Invalid offset")
				return
			}

			res, err := queries.GetPictures(r.Context(), db_sql.GetPicturesParams{
				Limit:  int64(limitAsInt),
				Offset: int64(offsetAsInt),
			})
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}

			count, err := queries.AllPicturesCount(r.Context())
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}

			resp := response{
				Total: int(count),
				Pictures: lo.Map(res, func(row db_sql.Picture, index int) picture {
					return picture{
						ID:        row.ID,
						BlobId:    row.BlobStorageID.String,
						YoutubeId: row.YtVideoID.String,
					}
				}),
			}

			render.JSON(w, r, resp)
		},
	)
}
