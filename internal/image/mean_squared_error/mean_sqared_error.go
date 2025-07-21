package mse

import (
	"math"

	"github.com/go-errors/errors"

	"gocv.io/x/gocv"
)

func GetMeanSquaredError(img1, img2 *gocv.Mat) (float64, error) {
	if img1.Rows() != img2.Rows() || img1.Cols() != img2.Cols() {
		return 0, errors.New("images do not have the same size")
	}

	converted16bit1 := gocv.NewMat()
	img1.ConvertToWithParams(&converted16bit1, gocv.MatTypeCV16UC3, 1, 0)
	defer converted16bit1.Close()

	converted16bit2 := gocv.NewMat()
	img2.ConvertToWithParams(&converted16bit2, gocv.MatTypeCV16UC3, 1, 0)
	defer converted16bit2.Close()

	diff := gocv.NewMat()
	defer diff.Close()
	gocv.AbsDiff(converted16bit1, converted16bit2, &diff)

	multiplied := gocv.NewMat()
	defer multiplied.Close()
	gocv.Multiply(diff, diff, &multiplied)

	sumVertical := gocv.NewMatFromScalar(gocv.NewScalar(0, 0, 0, 0), gocv.MatTypeCV64FC3)
	defer sumVertical.Close()
	gocv.Reduce(multiplied, &sumVertical, 0, gocv.ReduceSum, gocv.MatTypeCV64FC3)

	sumToScalar := gocv.NewMatFromScalar(gocv.NewScalar(0, 0, 0, 0), gocv.MatTypeCV64FC3)
	defer sumToScalar.Close()
	gocv.Reduce(sumVertical, &sumToScalar, 1, gocv.ReduceSum, gocv.MatTypeCV64FC3)
	sumScalar := sumToScalar.GetVecdAt(0, 0)

	mse := 0.0
	// Range over the channels
	for _, f := range sumScalar {
		mse += f
	}

	mse /= float64(img1.Rows() * img1.Cols() * img1.Channels())
	mse = math.Sqrt(mse) / 255.0

	return mse, nil
}
