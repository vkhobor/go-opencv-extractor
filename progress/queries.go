package progress

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/db"
)

type Queries struct {
	Queries *db.Queries
}

type VideoProgress struct {
	VideoId          string
	ImportProgress   int
	DownloadProgress int
	ImportError      string
	DownloadError    string
}

func (jc *Queries) VideoProgresses(ctx context.Context, jobId string) ([]VideoProgress, error) {
	res, err := jc.Queries.GetJobVideosWithProgress(ctx, sql.NullString{
		String: jobId,
		Valid:  true,
	})
	slog.Debug("res", "res", res, "err", err)
	if err != nil {
		return nil, err
	}

	uniqs := lo.UniqBy(res, func(item db.GetJobVideosWithProgressRow) string {
		return item.ID
	})
	videoIds := lo.Map(uniqs, func(item db.GetJobVideosWithProgressRow, i int) string {
		return item.ID
	})
	slog.Debug("videoIds", "videoIds", videoIds)
	result := make([]VideoProgress, 0, len(videoIds))

	for _, videoId := range videoIds {
		var importProgress, downloadProgress int
		var importError, downloadError string

		for _, v := range res {
			if v.ID == videoId {

				importAttempts := lo.FilterMap(res, func(item db.GetJobVideosWithProgressRow, i int) (Attempt, bool) {
					return Attempt{
						Progress: int(item.Progress_2.Int64),
						Error:    item.Error_2.String,
					}, item.ID_3.Valid
				})
				downloadAttempts := lo.FilterMap(res, func(item db.GetJobVideosWithProgressRow, i int) (Attempt, bool) {
					return Attempt{
						Progress: int(item.Progress.Int64),
						Error:    item.Error.String,
					}, item.ID_2.Valid
				})

				importProgress, importError = getProgress(importAttempts)
				downloadProgress, downloadError = getProgress(downloadAttempts)
			}
		}

		result = append(result, VideoProgress{
			VideoId:          videoId,
			ImportProgress:   importProgress,
			DownloadProgress: downloadProgress,
			ImportError:      importError,
			DownloadError:    downloadError,
		})
	}
	return result, nil
}

type Attempt struct {
	Progress int
	Error    string
}

func getProgress(v []Attempt) (int, string) {
	if len(v) == 0 {
		return 0, ""
	}

	current, _, ok := lo.FindIndexOf(v, func(item Attempt) bool {
		return item.Error == ""
	})

	if !ok {
		errored := lo.Filter(v, func(item Attempt, intex int) bool {
			return item.Error != ""
		})[0]
		return 0, errored.Error
	}

	return current.Progress, ""
}
