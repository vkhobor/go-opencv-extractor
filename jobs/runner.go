package jobs

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/vkhobor/go-opencv/db_sql"
)

type JobCreator struct {
	Scrape   func(args ScrapeArgs) <-chan ScrapedVideo
	VImport  func(...DownlodedVideo) <-chan ImportedVideo
	Download func(...ScrapedVideo) <-chan DownlodedVideo
	Queries  *db_sql.Queries
}

type ScrapeArgs struct {
	SearchQuery string
	Limit       int
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

// func (jc *JobCreator) CreateJob(searchQuery string, limit int, jobId string) {
// 	progress := jc.RunJob(searchQuery, limit, jobId)
// 	for p := range progress {
// 		jc.OnProgressUpdate(p)
// 	}
// }

func (jc *JobCreator) RunJobPool() {
	go jc.RunDownloadJob()
	go jc.RunImportJob()
}

func (jc *JobCreator) RunScrapeJob(searchQuery string, limit int, jobId string) {
	scrapeChan := jc.Scrape(ScrapeArgs{SearchQuery: searchQuery, Limit: limit})

	for video := range scrapeChan {
		jc.SaveSraped(video, jobId)
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
		result[i] = DownlodedVideo{ScrapedVideo: ScrapedVideo{ID: v.ID}, SavePath: v.Path.String}

		return result
	}
	return result
}

var statuses = []string{"errored", "completed"}

// func (jc *JobCreator) RunJob(searchQuery string, limit int, jobId string) <-chan Progress {
// 	progress := make(chan Progress)

// 	go func() {

// 		scrapeChan := jc.Scrape(ScrapeArgs{SearchQuery: searchQuery, Limit: limit})
// 		scrapeChanMap := MultiplexChan(scrapeChan, func(input ScrapedVideo) string {
// 			ok := jc.SaveSraped(input, jobId)
// 			if ok {
// 				return "completed"
// 			} else {
// 				return "errored"
// 			}
// 		}, statuses)
// 		scrapeChan = scrapeChanMap["completed"]

// 		downloadChan := jc.Download(scrapeChan)
// 		downloadChanMap := MultiplexChan(downloadChan, func(input DownlodedVideo) string {
// 			ok := jc.DownloadSaved(input)
// 			fmt.Printf("downloaded video %v\n", ok)
// 			if ok {
// 				return "completed"
// 			} else {
// 				return "errored"
// 			}
// 		}, statuses)
// 		downloadChan = downloadChanMap["completed"]

// 		importChan := jc.VImport(downloadChan)
// 		importChanMap := MultiplexChan(importChan, func(input ImportedVideo) string {
// 			ok := jc.SaveImported(input)
// 			if ok {
// 				return "completed"
// 			} else {
// 				return "errored"
// 			}
// 		}, statuses)
// 		importChan = importChanMap["completed"]
// 		importCountChan := CountChan(importChan)

// 		failedChan := MergeChans(CountChan(importChanMap["errored"]), CountChan(downloadChanMap["errored"]), CountChan(scrapeChanMap["errored"]))

// 		latest := LatestFromChans(importCountChan, failedChan)

// 		for l := range latest {
// 			progress <- Progress{
// 				Total:     limit,
// 				Errored:   l[1],
// 				Completed: l[0],
// 			}
// 		}

// 		close(progress)
// 	}()

// 	return progress
// }

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
		ID: blobId,
		Path: sql.NullString{
			String: video.SavePath,
			Valid:  true,
		},
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
			ID: blobID,
			Path: sql.NullString{
				String: path.Join("~/test", frame.Path),
				Valid:  true,
			},
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

func (jc *JobCreator) OnProgressUpdate(progress Progress) {

}

func RemoveAllPaths(files ...string) {
	for _, file := range files {
		_ = os.Remove(file)
	}
}
