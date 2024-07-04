package queries

type Frame struct {
	FrameNumber int
	Path        string
}

type ImportedVideo struct {
	DownlodedVideo
	ExtractedFrames []Frame
}

type DownlodedVideo struct {
	ScrapedVideo
	SavePath string
	Error    error
}

type Job struct {
	// TODO remove limit
	Limit       int
	JobID       string
	SearchQuery string
	FilterID    string
}

type ScrapedVideo struct {
	Job
	ID string
}
