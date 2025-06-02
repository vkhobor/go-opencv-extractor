package surf

import (
	"errors"

	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/image"
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

func newSURFImageMatcher(refs []gocv.Mat, options ...SURFImageMatcherOption) (*SURFImageMatcher, error) {
	defer func() {
		for _, img := range refs {
			img.Close()
		}
	}()

	surfAlgorithm := contrib.NewSURF()

	descriptors, err := getDescriptorsFromImages(surfAlgorithm, refs)
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

func NewSURFImageMatcher(refs []string, options ...SURFImageMatcherOption) (*SURFImageMatcher, error) {
	refImages, err := image.GetImagesFromPaths(refs)
	if err != nil {
		return nil, err
	}

	return newSURFImageMatcher(refImages, options...)
}

func NewSURFImageMatcherFromMats(refs []gocv.Mat, options ...SURFImageMatcherOption) (*SURFImageMatcher, error) {
	return newSURFImageMatcher(refs, options...)
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
	defer whiteMask.Close()

	_, descriptorsFrame := e.surfAlgorithm.DetectAndCompute(frameInGrayscale, whiteMask)
	defer descriptorsFrame.Close()

	// Each keypoint of the first image is matched with
	// a number of keypoints from the second image.
	// We keep the 2 best matches for each keypoint
	// (best matches = the ones with the smallest distance measurement).
	// Lowe's test checks that the two distances are sufficiently different.
	// If they are not, then the keypoint is eliminated and
	// will not be used for further calculations.
	knnMatches := getKnnMatches(e.matcher, e.descriptors, descriptorsFrame, 2)
	// David Lowe proposed a simple method for filtering keypoint matches by eliminating matches when the second-best match is almost as good.
	goodMatches := filterByDawidLoweRatioTest(knnMatches, e.ratioTestThreshold)

	everyHasSufficient := lo.EveryBy(goodMatches, func(item []gocv.DMatch) bool {
		return hasSufficientGoodMatches(item, e.minThresholdForSURFMatches, e.minSURFMatches)
	})

	return everyHasSufficient
}

func getDescriptorsFromImages(surf contrib.SURF, images []gocv.Mat) ([]gocv.Mat, error) {
	var descriptors []gocv.Mat
	for _, img := range images {
		mask := gocv.NewMat()
		defer mask.Close()
		_, descriptor := surf.DetectAndCompute(img, mask)
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

// Gets the best k matches from the query descriptors
func getKnnMatches(matcher gocv.BFMatcher, queryDescriptors []gocv.Mat, descriptorsFrame gocv.Mat, k int) [][][]gocv.DMatch {
	var knnMatches [][][]gocv.DMatch
	for _, descriptor := range queryDescriptors {
		knnMatches = append(knnMatches, matcher.KnnMatch(descriptorsFrame, descriptor, k))
	}
	return knnMatches
}

func filterByDawidLoweRatioTest(knnMatches [][][]gocv.DMatch, threshold float64) [][]gocv.DMatch {
	var goodMatches [][]gocv.DMatch
	for _, matches := range knnMatches {
		var goodMatch []gocv.DMatch
		for _, m := range matches {
			// Distances are sorted in ascending order, so the first match is the closest
			// 0 means identical, 1 means completely different
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
