package background

import (
	"github.com/vkhobor/go-opencv/features"
	"github.com/vkhobor/go-opencv/mlog"
)

func (d *DbMonitor) StartImport() {
	for video := range d.ImportInput {
		importer := features.ImportVideoFeature{
			Queries: d.Queries,
			Config:  d.Config,
		}

		mlog.Log().Debug("Importer starting importing", "video", video, "method", "Start")
		err := importer.ImportVideo(video.ID, video.JobID, video.FilterID)
		if err != nil {
			mlog.Log().Error("Error while importing video", "error", err, "video", video)
			continue
		}
	}
}
