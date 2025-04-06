package youtube

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/google/uuid"
	"github.com/iawia002/lux/downloader"
	"github.com/iawia002/lux/extractors"

	_ "github.com/iawia002/lux/extractors/youtube"
)

type YoutubeVideo string

func (y YoutubeVideo) String() string {
	return string(y)
}

func (y YoutubeVideo) URL() string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", y.String())
}

var youtubeRegexp = regexp.MustCompile(`^.*((youtu.be\/)|(v\/)|(\/u\/\w\/)|(embed\/)|(watch\?))\??v?=?([^#&?]*).*`)

func NewYoutubeIDFromUrl(url string) (YoutubeVideo, error) {
	match := youtubeRegexp.FindStringSubmatch(url)
	if len(match) > 7 && len(match[7]) == 11 {
		return YoutubeVideo(match[7]), nil
	}

	return "", errors.New("cannot parse url")
}

type youtubeClient string

const (
	AndroidClient  = youtubeClient("android")
	WebClient      = youtubeClient("web")
	EmbeddedClient = youtubeClient("embedded")
)

// Setting clients is based on a global variable, which is not thread-safe.
var m sync.Mutex

func (y YoutubeVideo) DownloadToFolder(clientType youtubeClient, folderPath string, progress chan<- float64) (string, error) {
	// m.Lock()

	// switch clientType {
	// case AndroidClient:
	// 	youtube.DefaultClient = youtube.AndroidClient
	// case WebClient:
	// 	youtube.DefaultClient = youtube.WebClient
	// case EmbeddedClient:
	// 	youtube.DefaultClient = youtube.EmbeddedClient
	// default:
	// 	youtube.DefaultClient = youtube.AndroidClient
	// }

	// client := youtube.Client{}
	// video, err := client.GetVideo(y.String())

	// m.Unlock()

	// if err != nil {
	// 	return "", err
	// }

	// video.FilterQuality("720p")

	// if len(video.Formats) == 0 {
	// 	return "", errors.New("no matching formats found")
	// }

	// stream, size, err := client.GetStream(video, &video.Formats[0])
	// if err != nil {
	// 	return "", err
	// }
	// defer stream.Close()

	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return "", err
	}
	id := uuid.New()
	fileName := fmt.Sprintf("%v_%v.mp4", id.String(), y.String())
	filePath := filepath.Join(folderPath, fileName)

	err = download(context.Background(), y.URL(), folderPath, fileName)
	if err != nil {
		return "", err
	}

	// file, err := os.Create(filePath)
	// if err != nil {
	// 	return "", err
	// }
	// defer file.Close()

	// progressReporter := &writeReporter{
	// 	Total:    size,
	// 	Progress: progress,
	// }
	// _, err = io.Copy(io.MultiWriter(file, progressReporter), stream)
	// if err != nil {
	// 	return "", err
	// }

	return filePath, nil
}

type writeReporter struct {
	current  int64
	Total    int64
	Progress chan<- float64
}

func (w *writeReporter) Write(p []byte) (n int, err error) {
	n = len(p)
	w.current += int64(n)

	w.Progress <- float64(w.current) / float64(w.Total)

	return
}

func download(c context.Context, videoURL string, output string, name string) error {
	data, err := extractors.Extract(videoURL, extractors.Options{})
	if err != nil {
		// if this error occurs, it means that an error occurred before actually starting to extract data
		// (there is an error in the preparation step), and the data list is empty.
		return err
	}

	defaultDownloader := downloader.New(downloader.Options{
		OutputPath: output,
		OutputName: name,
	})
	errors := make([]error, 0)
	for _, item := range data {
		if item.Err != nil {
			// if this error occurs, the preparation step is normal, but the data extraction is wrong.
			// the data is an empty struct.
			errors = append(errors, item.Err)
			continue
		}
		if err = defaultDownloader.Download(item); err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) != 0 {
		return errors[0]
	}
	return nil
}
