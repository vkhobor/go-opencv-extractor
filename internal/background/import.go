package background

import (
	"context"

	"github.com/vkhobor/go-opencv/internal/features/import_video"
	"github.com/vkhobor/go-opencv/internal/mlog"
)

func (d *DbMonitor) StartImport() {
	adapter := NewDbAdapter(d.SqlDB)
	for video := range d.ImportInput {
		importer := import_video.ImportVideoFeature{
			SqlDB:   adapter.TxEr,
			Querier: adapter.Querier,
			Config:  d.Config,
		}

		mlog.Log().Debug("Importer starting importing", "video", video, "method", "Start")
		err := importer.ImportVideo(context.Background(), video.ID, video.JobID, video.FilterID)
		if err != nil {
			mlog.Log().Error("Error while importing video", "error", err, "video", video)
			continue
		}
	}
}
