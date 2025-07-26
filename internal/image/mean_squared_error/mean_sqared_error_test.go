package mse

import (
	"math"
	"testing"

	"gocv.io/x/gocv"
)

func TestGetMeanSquaredError(t *testing.T) {
	blackImg := gocv.NewMatWithSize(100, 100, gocv.MatTypeCV8UC3)
	defer blackImg.Close()
	blackImg.SetTo(gocv.NewScalar(0, 0, 0, 0))

	whiteImg := gocv.NewMatWithSize(100, 100, gocv.MatTypeCV8UC3)
	defer whiteImg.Close()
	whiteImg.SetTo(gocv.NewScalar(255, 255, 255, 0))

	blackImg2 := gocv.NewMatWithSize(100, 100, gocv.MatTypeCV8UC3)
	defer blackImg2.Close()
	blackImg2.SetTo(gocv.NewScalar(0, 0, 0, 0))

	sizeDiffImage := gocv.NewMatWithSize(50, 50, gocv.MatTypeCV8UC3) // Different size
	defer sizeDiffImage.Close()

	tests := []struct {
		name      string
		img1      *gocv.Mat
		img2      *gocv.Mat
		wantErr   bool
		wantMSE   float64
		threshold float64
	}{
		{"Identical images", &blackImg, &blackImg2, false, 0.0, 0.0001},
		{"Different images", &blackImg, &whiteImg, false, 1, 0.0001},
		{"Size mismatch", &blackImg, &sizeDiffImage, true, 0.0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mse, err := GetMeanSquaredError(tt.img1, tt.img2)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetMeanSquaredError() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && math.Abs(mse-tt.wantMSE) > tt.threshold {
				t.Errorf("GetMeanSquaredError() mse = %v, wantMSE %v (threshold %v)", mse, tt.wantMSE, tt.threshold)
			}
		})
	}
}
