package domain

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
