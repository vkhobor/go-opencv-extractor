package surf

import (
	"errors"

	"github.com/samber/lo"
	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
)

type SURFImageMatcherConfig struct {
	ratioTestThreshold         float64
	minThresholdForSURFMatches float64
	minSURFMatches             int
}

var defaultCheckerConfig = SURFImageMatcherConfig{
	ratioTestThreshold:         0.5,
	minThresholdForSURFMatches: 0.3,
	minSURFMatches:             3,
}

func WithRatioThreshold(ratio float64) SURFImageMatcherOption {
	return func(c *SURFImageMatcher) {
		c.ratioTestThreshold = ratio
	}
}

func WithMinThreshold(threshold float64) SURFImageMatcherOption {
	return func(c *SURFImageMatcher) {
		c.minThresholdForSURFMatches = threshold
	}
}

func WithMinMatches(numOfMatches int) SURFImageMatcherOption {
	return func(c *SURFImageMatcher) {
		c.minSURFMatches = numOfMatches
	}
}

type SURFImageMatcherOption func(*SURFImageMatcher)

type SURFImageMatcher struct {
	SURFImageMatcherConfig
	descriptors   []gocv.Mat
	matcher       gocv.BFMatcher
	surfAlgorithm contrib.SURF
	Close         func() error
}

func NewSURFImageMatcher(refs []string, options ...SURFImageMatcherOption) (*SURFImageMatcher, error) {
	refImages, err := getImagesFromPaths(refs)
	if err != nil {
		return nil, err
	}
	defer func() {
		for _, img := range refImages {
			img.Close()
		}
	}()

	surfAlgorithm := contrib.NewSURF()

	descriptors, err := getDescriptorsFromImages(surfAlgorithm, refImages)
	if err != nil {
		return nil, err
	}

	matcher := gocv.NewBFMatcher()

	checker := &SURFImageMatcher{
		SURFImageMatcherConfig: defaultCheckerConfig,
		descriptors:            descriptors,
		matcher:                matcher,
		surfAlgorithm:          surfAlgorithm,
		Close: func() error {
			surfAlgorithm.Close()
			matcher.Close()

			for _, desc := range descriptors {
				err := desc.Close()
				if err != nil {
					return err
				}
			}
			return nil
		}}

	for _, option := range options {
		option(checker)
	}

	return checker, nil
}

func (e *SURFImageMatcher) IsImageMatch(frame *gocv.Mat) bool {
	frameInGrayscale := gocv.NewMat()
	defer frameInGrayscale.Close()

	if frame.Channels() == 3 {
		gocv.CvtColor(*frame, &frameInGrayscale, gocv.ColorBGRToGray)
	} else {
		frame.CopyTo(&frameInGrayscale)
	}

	whiteMask := fullWhiteMaskInSize(&frameInGrayscale)

	_, descriptorsFrame := e.surfAlgorithm.DetectAndCompute(frameInGrayscale, whiteMask)
	defer descriptorsFrame.Close()

	knnMatches := getKnnMatches(e.matcher, e.descriptors, descriptorsFrame)
	goodMatches := filterByDawidLoweRatioTest(knnMatches, e.ratioTestThreshold)

	everyHasSufficient := lo.EveryBy(goodMatches, func(item []gocv.DMatch) bool {
		return hasSufficientGoodMatches(item, e.minThresholdForSURFMatches, e.minSURFMatches)
	})

	return everyHasSufficient
}

func getImagesFromPaths(paths []string) ([]gocv.Mat, error) {
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

func getDescriptorsFromImages(surf contrib.SURF, images []gocv.Mat) ([]gocv.Mat, error) {
	var descriptors []gocv.Mat
	for _, img := range images {
		_, descriptor := surf.DetectAndCompute(img, gocv.NewMat())
		if descriptor.Empty() {
			for _, img := range descriptors {
				img.Close()
			}
			return nil, errors.New("Error computing descriptor")
		}
		descriptors = append(descriptors, descriptor)
	}
	return descriptors, nil
}

func getKnnMatches(matcher gocv.BFMatcher, descriptors []gocv.Mat, descriptorsFrame gocv.Mat) [][][]gocv.DMatch {
	var knnMatches [][][]gocv.DMatch
	for _, descriptor := range descriptors {
		knnMatches = append(knnMatches, matcher.KnnMatch(descriptorsFrame, descriptor, 2))
	}
	return knnMatches
}

func filterByDawidLoweRatioTest(knnMatches [][][]gocv.DMatch, threshold float64) [][]gocv.DMatch {
	var goodMatches [][]gocv.DMatch
	for _, matches := range knnMatches {
		var goodMatch []gocv.DMatch
		for _, m := range matches {
			if len(m) == 2 && m[0].Distance < threshold*m[1].Distance {
				goodMatch = append(goodMatch, m[0])
			}
		}
		goodMatches = append(goodMatches, goodMatch)
	}
	return goodMatches
}

func hasSufficientGoodMatches(goodMatches1 []gocv.DMatch, threshold float64, minMatches int) bool {
	return countMatchesBelowThreshold(goodMatches1, threshold) >= minMatches
}

func countMatchesBelowThreshold(matches []gocv.DMatch, threshold float64) int {
	count := 0
	for _, match := range matches {
		if match.Distance < threshold {
			count++
		}
	}
	return count
}

func fullWhiteMaskInSize(mat *gocv.Mat) gocv.Mat {
	mask := gocv.NewMatWithSizesWithScalar([]int{mat.Rows(), mat.Cols()}, gocv.MatTypeCV8U, gocv.NewScalar(255, 255, 255, 0))
	return mask
}
