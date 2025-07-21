package image

import (
	"errors"
	"io"

	"gocv.io/x/gocv"
)

func ReadImageFromReader(reader io.Reader) (gocv.Mat, error) {
	imgBytes, err := io.ReadAll(reader)
	if err != nil {
		return gocv.Mat{}, err
	}

	return gocv.IMDecode(imgBytes, gocv.IMReadColor)
}

func GetImagesFromPaths(paths []string) ([]gocv.Mat, error) {
	var images []gocv.Mat
	for _, path := range paths {
		img := gocv.IMRead(path, gocv.IMReadColor)
		if img.Empty() {
			for _, img := range images {
				img.Close()
			}
			return nil, errors.New("Error reading image")
		}
		images = append(images, img)
	}
	return images, nil
}
