package features

import (
	"github.com/vkhobor/go-opencv/queries"
	"github.com/vkhobor/go-opencv/youtube"
)

type ScrapeByIdFeature struct {
	Queries              *queries.Queries
	MaxErrorStopRetrying int
}

func (i *ScrapeByIdFeature) ScrapeByIdSearch(id youtube.YoutubeVideo, jobID string, limit int) error {
	// TODO check if video is already scraped, optionally abort while progressing

	// TODO make this more efficient, no need to query for every video
	// jobs := i.Queries.GetToScrapeVideos()
	// actualJob := queries.Job{}
	// for _, job := range jobs {
	// 	if job.JobID == jobID {
	// 		actualJob = job
	// 		break
	// 	}
	// }

	// err := i.Queries.SaveNewlyScraped(jobID, string(id), actualJob.FilterID)
	// if err != nil {
	// 	return err
	// }

	return nil
}
