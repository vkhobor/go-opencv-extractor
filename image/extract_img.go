package image

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math"
	"os"

	"github.com/vkhobor/go-opencv/memo"
	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
)

type Config struct {
	VideoPath   string
	OriginalFPS float64
	WantFPS     int
}

//go:embed horiz.png
var horiz []byte

//go:embed vertic.png
var vertic []byte

type ExtractIterator struct {
	descriptorRef1    gocv.Mat
	descriptorRef2    gocv.Mat
	memoizedallowMask func(orig *gocv.Mat) gocv.Mat
	modFPS            int
	matcher           gocv.BFMatcher
	videoCapture      *gocv.VideoCapture
	currentFrame      gocv.Mat
	surf              contrib.SURF
	cfg               Config
	progressChan      chan<- int
	Close             func()
}

func NewExtractIterator(cfg Config, progressChan chan<- int) (*ExtractIterator, error) {
	imageReader := bytes.NewReader(vertic)

	// Decode the image
	img1, _, err := image.Decode(imageReader)
	if err != nil {
		log.Fatal(err)
	}

	imageReader2 := bytes.NewReader(horiz)

	// Decode the image
	img2, _, err := image.Decode(imageReader2)
	if err != nil {
		log.Fatal(err)
	}

	refImage1, _ := gocv.ImageToMatRGBA(img1)
	refImage2, _ := gocv.ImageToMatRGBA(img2)

	surf := contrib.NewSURF()
	matcher := gocv.NewBFMatcher()
	_, descriptorsRef1 := surf.DetectAndCompute(refImage1, gocv.NewMat())
	_, descriptorsRef2 := surf.DetectAndCompute(refImage2, gocv.NewMat())

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
		surf:              surf,
		matcher:           matcher,
		currentFrame:      frame,
		descriptorRef1:    descriptorsRef1,
		descriptorRef2:    descriptorsRef2,
		memoizedallowMask: memoizedAllowMask,
		progressChan:      progressChan,
		cfg:               cfg,
		Close: func() {
			closeMemo()
			frame.Close()
			surf.Close()
			matcher.Close()
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
		e.progressChan <- e.CurrentFrame()

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

func (e *ExtractIterator) checkFrame() bool {
	frame := e.currentFrame
	surf := e.surf
	matcher := e.matcher
	descriptorsRef1 := e.descriptorRef1
	descriptorsRef2 := e.descriptorRef2
	allowMask := e.memoizedallowMask

	grayFrame := gocv.NewMat()
	defer grayFrame.Close()
	if frame.Channels() == 3 {
		gocv.CvtColor(frame, &grayFrame, gocv.ColorBGRToGray)
	} else {
		frame.CopyTo(&grayFrame)
	}

	noneMask := allowMask(&grayFrame)

	_, descriptorsFrame := surf.DetectAndCompute(grayFrame, noneMask)
	defer descriptorsFrame.Close()

	knnMatches1 := matcher.KnnMatch(descriptorsFrame, descriptorsRef1, 2)
	knnMatches2 := matcher.KnnMatch(descriptorsFrame, descriptorsRef2, 2)

	const ratioThreshold = 0.5
	var goodMatches1, goodMatches2 []gocv.DMatch
	for _, m := range knnMatches1 {
		if len(m) == 2 && m[0].Distance < ratioThreshold*m[1].Distance {
			goodMatches1 = append(goodMatches1, m[0])
		}
	}
	for _, m := range knnMatches2 {
		if len(m) == 2 && m[0].Distance < ratioThreshold*m[1].Distance {
			goodMatches2 = append(goodMatches2, m[0])
		}
	}

	if !hasSufficientGoodMatches(goodMatches1, goodMatches2, 0.3, 3) {
		return false
	}

	return true
}

// hasSufficientGoodMatches checks if there are enough good matches in both sets.
func hasSufficientGoodMatches(goodMatches1, goodMatches2 []gocv.DMatch, threshold float64, minMatches int) bool {
	return countMatchesBelowThreshold(goodMatches1, threshold) >= minMatches &&
		countMatchesBelowThreshold(goodMatches2, threshold) >= minMatches
}

// countMatchesBelowThreshold counts the number of matches below a specified distance threshold.
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
