package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/db_sql"
)

func HandleGetStats(queries *db_sql.Queries) http.HandlerFunc {
	type statsResponse struct {
		VideosChecked         int `json:"videos_checked"`
		MatchingPicturesSaved int `json:"macthing_pictures_saved"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			res, err := queries.ListImportedVideosWithSaved(r.Context())
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}

			sumFrames := lo.Reduce(res, func(acc int, v db_sql.ListImportedVideosWithSavedRow, _ int) int {
				return acc + int(v.Importedpictures)
			}, 0)

			resp := statsResponse{
				VideosChecked:         len(res),
				MatchingPicturesSaved: sumFrames,
			}

			render.JSON(w, r, resp)
		},
	)
}
