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
	SavePath       string
	Error          error
	ImportProgress int64
}

type Job struct {
	Limit       int
	JobID       string
	SearchQuery string
	YouTubeID   string
	Name        string
	FilterID    string
}

type ScrapedVideo struct {
	Job
	ID string
}
