package youtube

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/google/uuid"
	"github.com/kkdai/youtube/v2"
)

var youtubeRegexp = regexp.MustCompile(`^.*((youtu.be\/)|(v\/)|(\/u\/\w\/)|(embed\/)|(watch\?))\??v?=?([^#&?]*).*`)

type YoutubeVideo string

func (y YoutubeVideo) String() string {
	return string(y)
}

func NewYoutubeIDFromUrl(url string) (YoutubeVideo, error) {
	match := youtubeRegexp.FindStringSubmatch(url)
	if len(match) > 7 && len(match[7]) == 11 {
		return YoutubeVideo(match[7]), nil
	}

	return "", errors.New("cannot parse url")
}

func (y YoutubeVideo) DownloadToFolder(folderPath string, progress chan<- float64) (string, error) {
	youtube.DefaultClient = youtube.WebClient
	client := youtube.Client{}

	video, err := client.GetVideo(y.String())
	if err != nil {
		return "", err
	}

	video.FilterQuality("720p")

	if len(video.Formats) == 0 {
		return "", errors.New("no matching formats found")
	}

	stream, size, err := client.GetStream(video, &video.Formats[0])
	if err != nil {
		return "", err
	}
	defer stream.Close()

	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return "", err
	}
	id := uuid.New()
	fileName := fmt.Sprintf("%v_%v.mp4", id.String(), y.String())
	filePath := filepath.Join(folderPath, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	progressReporter := &WriteReporter{
		Total:    size,
		Progress: progress,
	}
	_, err = io.Copy(io.MultiWriter(file, progressReporter), stream)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

type WriteReporter struct {
	Total    int64
	current  int64
	Progress chan<- float64
}

func (w *WriteReporter) Write(p []byte) (n int, err error) {
	n = len(p)
	w.current += int64(n)

	select {
	case w.Progress <- float64(w.current) / float64(w.Total) * 100:
	default:
	}

	return
}
