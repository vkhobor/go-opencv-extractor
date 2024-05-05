package jobs

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/db_sql"
)

type JobCreator struct {
	Scrape   func(args ScrapeArgs, ctx context.Context) <-chan ScrapedVideo
	VImport  func(refImagePaths []string, videos ...DownlodedVideo) <-chan ImportedVideo
	Download func(...ScrapedVideo) <-chan DownlodedVideo
	Queries  *db_sql.Queries
}

type ScrapeArgs struct {
	Limit       int
	JobId       string
	SearchQuery string
}

type ScrapedVideo struct {
	ID string
}

type DownlodedVideo struct {
	ScrapedVideo
	SavePath string
	Error    error
}

type Frame struct {
	FrameNumber int
	Path        string
}

type ImportedVideo struct {
	DownlodedVideo
	ExtractedFrames []Frame
	Error           error
}

type Progress struct {
	Total     int
	Completed int
	Errored   int
}

func (jc *JobCreator) RunJobPoolOnce() {
	slog.Info("Running job pool")
	jc.RunScrapeJob()
	jc.RunDownloadJob()
	jc.RunImportJob()
}

func (jc *JobCreator) RunScrapeJob() {
	toScrape := jc.GetToScrapeVideos()
	if len(toScrape) == 0 {
		slog.Debug("No videos to scrape")
		return
	}

	slog.Info("Running scrape job", "needed_to_scrape", toScrape)
	for _, scrapeArgs := range toScrape {
		toFind := scrapeArgs.Limit

		if toFind <= 0 {
			continue
		}

		ctx, cancel := context.WithCancel(context.Background())
		scrapeChan := jc.Scrape(ScrapeArgs{SearchQuery: scrapeArgs.SearchQuery, Limit: scrapeArgs.Limit}, ctx)

		for video := range scrapeChan {
			if toFind <= 0 {
				continue
			} else {
				if jc.SaveSraped(video, scrapeArgs.JobId) {
					toFind--
					if toFind <= 0 {
						cancel()
					}
				}
			}

		}
		cancel()
	}
}

func (jc *JobCreator) RunDownloadJob() {
	scraped := jc.GetScrapedVideos()
	if len(scraped) == 0 {
		slog.Debug("No videos to download")
		return
	}

	slog.Info("Running download job")
	downloadChan := jc.Download(scraped...)

	for video := range downloadChan {
		jc.DownloadSaved(video)
	}
}

func (jc *JobCreator) RunImportJob() {

	d := jc.GetDownloadedVideos()
	refs, err := jc.GetRefImages()
	if err != nil || len(refs) == 0 || len(d) == 0 {
		slog.Debug("No videos to import or no reference images")
		return
	}

	slog.Info("Running import job", "videos_waiting_for_import", d)
	importChan := jc.VImport(refs, d...)

	for video := range importChan {
		jc.SaveImported(video)
	}
}

func (jc *JobCreator) GetToScrapeVideos() []ScrapeArgs {
	dbVal, err := jc.Queries.GetToScrapeVideos(context.Background())

	if err != nil {
		return []ScrapeArgs{}
	}

	return lo.FilterMap(dbVal, func(item db_sql.GetToScrapeVideosRow, i int) (ScrapeArgs, bool) {
		return ScrapeArgs{
			SearchQuery: item.SearchQuery.String,
			Limit:       int(item.Limit.Int64 - item.FoundVideos),
			JobId:       item.ID,
		}, item.Limit.Int64-item.FoundVideos > 0
	})
}

func (jc *JobCreator) GetScrapedVideos() []ScrapedVideo {
	val, err := jc.Queries.GetScrapedVideos(context.Background())
	if err != nil {
		return []ScrapedVideo{}
	}

	result := make([]ScrapedVideo, len(val))
	for i, v := range val {
		result[i] = ScrapedVideo{ID: v.ID}
	}

	return result
}

func (jc *JobCreator) GetDownloadedVideos() []DownlodedVideo {
	val, err := jc.Queries.GetVideosDownloaded(context.Background())
	if err != nil {
		return []DownlodedVideo{}
	}
	result := make([]DownlodedVideo, len(val))
	for i, v := range val {
		result[i] = DownlodedVideo{ScrapedVideo: ScrapedVideo{ID: v.ID}, SavePath: v.Path}
	}
	return result
}

func (jc *JobCreator) GetRefImages() ([]string, error) {
	val, err := jc.Queries.GetReferences(context.Background())
	if err != nil {
		return nil, err
	}

	return lo.Map(val, func(item db_sql.GetReferencesRow, i int) string {
		return item.Path
	}), nil
}

func (jc *JobCreator) SaveSraped(video ScrapedVideo, jobId string) bool {
	_, err := jc.Queries.AddYtVideo(context.Background(), db_sql.AddYtVideoParams{
		ID: video.ID,
		JobID: sql.NullString{
			String: jobId,
			Valid:  true,
		},
		Status: sql.NullString{
			String: "scraped",
			Valid:  true,
		},
	})

	if err != nil {
		return false
	}
	return true
}

func (jc *JobCreator) DownloadSaved(video DownlodedVideo) {
	slog.Debug("Saving downloaded", "video", video)
	if video.Error != nil {
		_, err := jc.Queries.UpdateStatus(context.Background(), db_sql.UpdateStatusParams{
			ID: video.ID,
			Status: sql.NullString{
				String: "errored",
				Valid:  true,
			},
			Error: sql.NullString{
				Valid:  true,
				String: video.Error.Error(),
			},
		})
		if err != nil {
			slog.Error("Error while updating status", "error", err)
		}

		return
	}
	blobId := uuid.New().String()
	_, errAddBlob := jc.Queries.AddBlob(context.Background(), db_sql.AddBlobParams{
		ID:   blobId,
		Path: video.SavePath,
	})

	_, errUpdateStatus := jc.Queries.UpdateStatus(context.Background(), db_sql.UpdateStatusParams{
		ID: video.ID,
		Status: sql.NullString{
			String: "downloaded",
			Valid:  true,
		},
	})

	jc.Queries.AddBlobToVideo(context.Background(), db_sql.AddBlobToVideoParams{
		BlobStorageID: sql.NullString{
			String: blobId,
			Valid:  true,
		},
		ID: video.ID,
	})

	if errAddBlob != nil || errUpdateStatus != nil {
		slog.Error("Error while updating status", "error", errAddBlob, "error2", errUpdateStatus)
		RemoveAllPaths(video.SavePath)
		return
	}
}

func (jc *JobCreator) SaveImported(video ImportedVideo) {
	if video.Error != nil {
		jc.Queries.UpdateStatus(context.Background(), db_sql.UpdateStatusParams{
			ID: video.ID,
			Status: sql.NullString{
				String: "errored",
				Valid:  true,
			},
		})
		return
	}

	paths := []string{}
	for _, frame := range video.ExtractedFrames {
		paths = append(paths, frame.Path)

		blobID := uuid.New().String()
		jc.Queries.AddBlob(context.Background(), db_sql.AddBlobParams{
			ID:   blobID,
			Path: path.Join("~/test", frame.Path),
		})
		jc.Queries.AddPicture(context.Background(), db_sql.AddPictureParams{
			ID: uuid.New().String(),
			YtVideoID: sql.NullString{
				String: video.ID,
				Valid:  true,
			},
			FrameNumber: sql.NullInt64{
				Int64: int64(frame.FrameNumber),
				Valid: true,
			},
			BlobStorageID: sql.NullString{
				String: blobID,
				Valid:  true,
			},
		})
	}

	jc.Queries.UpdateStatus(
		context.Background(),
		db_sql.UpdateStatusParams{
			Status: sql.NullString{
				String: "imported",
				Valid:  true,
			},
			ID: video.ID,
		})

	return
}

func RemoveAllPaths(files ...string) {
	for _, file := range files {
		_ = os.Remove(file)
	}
}
