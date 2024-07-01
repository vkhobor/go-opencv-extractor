package image

import (
	_ "embed"
	"errors"
	"fmt"
	_ "image/png"
	"math"
	"os"

	"github.com/samber/lo"
	"github.com/vkhobor/go-opencv/memo"
	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
)

type Config struct {
	VideoPath        string
	OriginalFPS      float64
	WantFPS          int
	PathsToRefImages []string
}

var (
	surfAlgorithm = contrib.NewSURF()
	matcher       = gocv.NewBFMatcher()
)

type ExtractIterator struct {
	descriptors       []gocv.Mat
	memoizedallowMask func(orig *gocv.Mat) gocv.Mat
	modFPS            int
	videoCapture      *gocv.VideoCapture
	currentFrame      gocv.Mat
	cfg               Config
	onProgress        func(*ExtractIterator)
	Close             func()
}

func GetImagesFromPaths(paths []string) ([]gocv.Mat, error) {
	var images []gocv.Mat
	for _, path := range paths {
		img := gocv.IMRead(path, gocv.IMReadColor)
		if img.Empty() {
			return nil, errors.New("Error reading image")
		}
		images = append(images, img)
	}
	return images, nil
}

func GetDescriptorsFromImages(images []gocv.Mat) ([]gocv.Mat, error) {
	var descriptors []gocv.Mat
	for _, img := range images {
		_, descriptor := surfAlgorithm.DetectAndCompute(img, gocv.NewMat())
		descriptors = append(descriptors, descriptor)
	}
	return descriptors, nil
}

func NewExtractIterator(cfg Config, onProgress func(*ExtractIterator)) (*ExtractIterator, error) {
	refImages, _ := GetImagesFromPaths(cfg.PathsToRefImages)
	GetDescriptorsFromImages(refImages)

	descriptors, _ := GetDescriptorsFromImages(refImages)

	// Check if the video file exists
	if _, err := os.Stat(cfg.VideoPath); os.IsNotExist(err) {
		return nil, errors.New("Video file does not exist")
	}

	capture, err := gocv.OpenVideoCapture(cfg.VideoPath)
	if err != nil {
		return nil, err
	}

	frame := gocv.NewMat()

	memoizedAllowMask, closeMemo := memo.Memoize(createFullMask, func(mat *gocv.Mat) string { return fmt.Sprint(mat.Cols()) + fmt.Sprint(mat.Rows()) })

	modFPS := int(math.Ceil(float64(cfg.OriginalFPS) / float64(cfg.WantFPS)))

	return &ExtractIterator{
		videoCapture:      capture,
		modFPS:            modFPS,
		descriptors:       descriptors,
		currentFrame:      frame,
		memoizedallowMask: memoizedAllowMask,
		onProgress:        onProgress,
		cfg:               cfg,
		Close: func() {
			closeMemo()
			frame.Close()
		},
	}, nil

}

func (e *ExtractIterator) Value() gocv.Mat {
	return e.currentFrame
}

func (e *ExtractIterator) Length() int {
	return int(e.videoCapture.Get(gocv.VideoCaptureFrameCount))
}

func (e *ExtractIterator) CurrentFrame() int {
	return int(e.videoCapture.Get(gocv.VideoCapturePosFrames))
}

func (e *ExtractIterator) Next() bool {
	for {
		success := e.videoCapture.Read(&e.currentFrame)
		if !success || e.currentFrame.Empty() {
			return false
		}
		e.onProgress(e)

		frameNum := int(e.videoCapture.Get(gocv.VideoCapturePosFrames))
		if frameNum%e.modFPS != 0 {
			continue
		}

		ok := e.checkFrame()

		if ok {
			return true
		}
	}
}

func GetKnnMatches(descriptors []gocv.Mat, descriptorsFrame gocv.Mat) [][][]gocv.DMatch {
	var knnMatches [][][]gocv.DMatch
	for _, descriptor := range descriptors {
		knnMatches = append(knnMatches, matcher.KnnMatch(descriptorsFrame, descriptor, 2))
	}
	return knnMatches
}

func FilterGoodMatches(knnMatches [][][]gocv.DMatch, threshold float64) [][]gocv.DMatch {
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

func (e *ExtractIterator) checkFrame() bool {
	frame := e.currentFrame
	allowMask := e.memoizedallowMask

	grayFrame := gocv.NewMat()
	defer grayFrame.Close()
	if frame.Channels() == 3 {
		gocv.CvtColor(frame, &grayFrame, gocv.ColorBGRToGray)
	} else {
		frame.CopyTo(&grayFrame)
	}

	noneMask := allowMask(&grayFrame)

	_, descriptorsFrame := surfAlgorithm.DetectAndCompute(grayFrame, noneMask)
	defer descriptorsFrame.Close()

	knnMatches := GetKnnMatches(e.descriptors, descriptorsFrame)
	const ratioThreshold = 0.5
	goodMatches := FilterGoodMatches(knnMatches, ratioThreshold)

	const minThreshold, minMatches = 0.3, 3
	everyHasSufficient := lo.EveryBy(goodMatches, func(item []gocv.DMatch) bool {
		return hasSufficientGoodMatches(item, minThreshold, minMatches)
	})

	return everyHasSufficient
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

func createFullMask(img *gocv.Mat) gocv.Mat {

	mask := gocv.NewMatWithSizesWithScalar([]int{img.Rows(), img.Cols()}, gocv.MatTypeCV8U, gocv.NewScalar(255, 255, 255, 0))
	return mask
}
