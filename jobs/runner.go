package jobs

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/db_sql"
)

type JobCreator struct {
	Scrape   func(args ScrapeArgs, ctx context.Context) <-chan ScrapedVideo
	VImport  func(...DownlodedVideo) <-chan ImportedVideo
	Download func(...ScrapedVideo) <-chan DownlodedVideo
	Queries  *db_sql.Queries
}

type ScrapeArgs struct {
	Limit       int
	JobId       string
	SearchQuery string
	Offset      int
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

func (jc *JobCreator) RunJobPool() {
	go jc.RunScrapeJob()
	go jc.RunDownloadJob()
	go jc.RunImportJob()
}

func (jc *JobCreator) RunScrapeJob() {
	for {
		time.Sleep(5 * time.Second)

		toScrape := jc.GetToScrapeVideos()
		fmt.Printf("Running scrape job\n, %v\n", toScrape)
		for _, scrapeArgs := range toScrape {
			toFind := scrapeArgs.Limit

			if toFind <= 0 {
				continue
			}

			ctx, cancel := context.WithCancel(context.Background())
			scrapeChan := jc.Scrape(ScrapeArgs{SearchQuery: scrapeArgs.SearchQuery, Limit: scrapeArgs.Limit, Offset: scrapeArgs.Offset}, ctx)

			for video := range scrapeChan {
				fmt.Println("Scraped video", toFind, video.ID)
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
		}
	}
}

func (jc *JobCreator) RunDownloadJob() {
	for {
		time.Sleep(5 * time.Second)

		fmt.Printf("Running download job\n")
		scraped := jc.GetScrapedVideos()
		downloadChan := jc.Download(scraped...)

		for video := range downloadChan {
			jc.DownloadSaved(video)
		}
	}
}

func (jc *JobCreator) RunImportJob() {
	for {
		time.Sleep(5 * time.Second)

		d := jc.GetDownloadedVideos()
		fmt.Printf("Running import job\n, %v\n", d)
		importChan := jc.VImport(d...)

		for video := range importChan {
			jc.SaveImported(video)
		}
	}
}

func (jc *JobCreator) GetToScrapeVideos() []ScrapeArgs {
	dbVal, err := jc.Queries.GetToScrapeVideos(context.Background())

	if err != nil {
		return []ScrapeArgs{}
	}

	return lo.Map(dbVal, func(item db_sql.GetToScrapeVideosRow, i int) ScrapeArgs {
		return ScrapeArgs{
			SearchQuery: item.SearchQuery.String,
			Limit:       int(item.Limit.Int64 - item.FoundVideos),
			Offset:      int(item.FoundVideos),
			JobId:       item.ID,
		}
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

		return result
	}
	return result
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
	blobId := uuid.New().String()
	_, err1 := jc.Queries.AddBlob(context.Background(), db_sql.AddBlobParams{
		ID:   blobId,
		Path: video.SavePath,
	})

	_, err2 := jc.Queries.UpdateStatus(context.Background(), db_sql.UpdateStatusParams{
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

	if err1 != nil || err2 != nil {
		RemoveAllPaths(video.SavePath)
		return
	}
	return
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
