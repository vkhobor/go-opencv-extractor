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
	"github.com/vkhobor/go-opencv/config"
)

func DownloadVideo(videoID string) (path string, title string, err error) {

	youtube.DefaultClient = youtube.AndroidClient
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		return "", "", err
	}

	video.FilterQuality("720p")

	if len(video.Formats) == 0 {
		return "", "", errors.New("no matching formats found")
	}

	stream, _, err := client.GetStream(video, &video.Formats[0])
	if err != nil {
		return "", "", err
	}
	defer stream.Close()

	err = os.MkdirAll(config.VideosDir, os.ModePerm)
	if err != nil {
		return "", "", err
	}
	id := uuid.New()
	fileName := fmt.Sprintf("%v_%v.mp4", id.String(), videoID)
	filePath := filepath.Join(config.VideosDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
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
