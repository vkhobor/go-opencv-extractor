package jobs

import (
	"log/slog"
	"time"
)

type JobManager struct {
	Wake           <-chan struct{}
	AutoWakePeriod time.Duration
	Scraper        Scraper
	Importer       *Importer
	Downloader     Downloader
}

func (jm *JobManager) Start() {

	ticker := time.NewTicker(jm.AutoWakePeriod)

	for {
		select {
		case <-jm.Wake:
			jm.RunPipelineOnce()
		case <-ticker.C:
			jm.RunPipelineOnce()
		}
	}
}

// RunPipelineOnce pulls all the jobs from the database and runs the pipeline on them once in order
func (jm *JobManager) RunPipelineOnce() {
	slog.Info("Running job pool")
	jm.Scraper()
	jm.Downloader()
	jm.Importer.ImportAllImportableFromDb()
}
