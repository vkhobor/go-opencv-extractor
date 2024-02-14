package importing

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/google/uuid"
	"github.com/kkdai/youtube/v2"
	"github.com/schollz/progressbar/v3"
)

func DownloadVideo(videoID string) (path string, title string, err error) {

	youtube.DefaultClient = youtube.AndroidClient
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		return "", "", err
	}

	video.FilterQuality("720p")

	stream, size, err := client.GetStream(video, &video.Formats[0])
	if err != nil {
		return "", "", err
	}
	defer stream.Close()

	tempDir := os.TempDir()
	id := uuid.New()
	fileName := fmt.Sprintf("%v_%v.mp4", id.String(), videoID)
	filePath := filepath.Join(tempDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	bar := progressbar.DefaultBytes(size)
	defer bar.Finish()

	_, err = io.Copy(io.MultiWriter(file, bar), stream)
	if err != nil {
		return "", "", err
	}
	return filePath, video.Title, nil
}

var youtubeRegexp = regexp.MustCompile(`^.*((youtu.be\/)|(v\/)|(\/u\/\w\/)|(embed\/)|(watch\?))\??v?=?([^#&?]*).*`)

func YoutubeParser(url string) (string, error) {
	match := youtubeRegexp.FindStringSubmatch(url)
	if len(match) > 7 && len(match[7]) == 11 {
		return match[7], nil
	}

	return "", errors.New("cannot parse url")
}
