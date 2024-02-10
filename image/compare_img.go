package image

import (
	"math"

	"github.com/go-errors/errors"

	"gocv.io/x/gocv"
)

func CompareImagesPath(path1, path2 string) (float64, error) {
	// Open the images
	img1 := gocv.IMRead(path1, gocv.IMReadColor)
	defer img1.Close()

	img2 := gocv.IMRead(path2, gocv.IMReadColor)
	defer img2.Close()

	return CompareImages(img1, img2)
}

func CompareImages(img1, img2 gocv.Mat) (float64, error) {

	// Check if the images have the same size
	if img1.Rows() != img2.Rows() || img1.Cols() != img2.Cols() {
		return 0, errors.New("images do not have the same size")
	}

	// Calculate the MSE
	var mse float64
	for y := 0; y < img1.Rows(); y++ {
		for x := 0; x < img1.Cols(); x++ {
			pixel1 := img1.GetVecbAt(y, x)
			pixel2 := img2.GetVecbAt(y, x)

			for i := 0; i < img1.Channels(); i++ {
				diff := float64(pixel1[i]) - float64(pixel2[i])
				mse += diff * diff
			}
		}
	}

	mse /= float64(img1.Rows() * img1.Cols() * img1.Channels())

	// Normalize the MSE to [0, 1] range
	mse = math.Sqrt(mse) / 255.0

	return mse, nil
}
