package surf

import (
	"fmt"
	"testing"

	"gocv.io/x/gocv"
)

func loadSampleImages() ([]gocv.Mat, error) {
	images := []string{
		"../../../samples/horiz.png",
		"../../../samples/vertic.png",
	}

	var mats []gocv.Mat
	for _, imgPath := range images {
		mat := gocv.IMRead(imgPath, gocv.IMReadGrayScale)
		if mat.Empty() {
			return nil, fmt.Errorf("failed to load image: %s", imgPath)
		}
		mats = append(mats, mat)
	}
	return mats, nil
}

func TestNewSURFImageMatcher(t *testing.T) {
	refs, err := loadSampleImages()
	if err != nil {
		t.Fatalf("failed to load sample images: %v", err)
	}
	defer func() {
		for _, mat := range refs {
			mat.Close()
		}
	}()

	matcher, err := NewSURFImageMatcherFromMats(refs)
	if err != nil {
		t.Fatalf("failed to create SURFImageMatcher: %v", err)
	}
	defer matcher.Close()

	if matcher == nil {
		t.Fatal("expected matcher to be non-nil")
	}
}

func loadTestImage(path string) (gocv.Mat, error) {
	mat := gocv.IMRead(path, gocv.IMReadGrayScale)
	if mat.Empty() {
		return gocv.Mat{}, fmt.Errorf("failed to load image: %s", path)
	}
	return mat, nil
}

func TestSURFImageMatcherShouldMatch(t *testing.T) {
	refs, err := loadSampleImages()
	if err != nil {
		t.Fatalf("failed to load sample images: %v", err)
	}
	defer func() {
		for _, mat := range refs {
			mat.Close()
		}
	}()

	matcher, err := NewSURFImageMatcherFromMats(refs)
	if err != nil {
		t.Fatalf("failed to create SURFImageMatcher: %v", err)
	}
	defer matcher.Close()

	testImages := []string{
		"../../../samples/should-match.jpg",
		"../../../samples/should-match2.jpg",
	}

	for _, imgPath := range testImages {
		testImg, err := loadTestImage(imgPath)
		if err != nil {
			t.Fatalf("failed to load test image: %v", err)
		}
		defer testImg.Close()

		if !matcher.IsImageMatch(&testImg) {
			t.Errorf("expected image %s to match, but it did not", imgPath)
		}
	}
}
