package video

import (
	"bytes"
	"errors"
	"io"
	"os"
	"sync"

	"gocv.io/x/gocv"
)

type FrameExtractor struct {
	// VideoCapture is not thread safe
	mu      sync.Mutex
	capture *gocv.VideoCapture
}

func (f *FrameExtractor) IsZero() bool {
	return f.capture == nil
}

func NewFrameExtractor(path string) (FrameExtractor, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return FrameExtractor{}, errors.New("Video file does not exist")
	}

	capture, err := gocv.OpenVideoCapture(path)
	if err != nil {
		return FrameExtractor{}, err
	}

	if capture.IsOpened() == false {
		return FrameExtractor{}, errors.New("Video could not be opened")
	}

	return FrameExtractor{
		capture: capture,
	}, nil
}

func (e *FrameExtractor) GetFrameAsMat(frame int) (gocv.Mat, error) {
	if e.IsZero() {
		return gocv.Mat{}, errors.New("FrameExtractor is zero")
	}
	if frame < 0 {
		return gocv.Mat{}, errors.New("Frame number must be non-negative")
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	e.capture.Set(gocv.VideoCapturePosFrames, float64(frame))

	mat := gocv.NewMat()

	ok := e.capture.Read(&mat)

	if !ok {
		return gocv.Mat{}, errors.New("Could not read frame")
	}

	return mat, nil
}

func (e *FrameExtractor) GetFrameAsJpeg(frame int) (io.ReadCloser, error) {
	if e.IsZero() {
		return nil, errors.New("FrameExtractor is zero")
	}
	if frame < 0 {
		return nil, errors.New("Frame number must be non-negative")
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	e.capture.Set(gocv.VideoCapturePosFrames, float64(frame))

	mat := gocv.NewMat()
	defer mat.Close()

	ok := e.capture.Read(&mat)

	if !ok {
		return nil, errors.New("Could not read frame")
	}

	return matToReadCloserInMemory(mat)
}

func matToReadCloserInMemory(mat gocv.Mat) (io.ReadCloser, error) {
	data, err := gocv.IMEncode(gocv.JPEGFileExt, mat)
	if err != nil {
		return nil, err
	}

	return newReadCloser(bytes.NewReader(data.GetBytes()), func() error {
		data.Close()
		return nil
	}), nil
}

func (e *FrameExtractor) Close() {
	if e.capture != nil {
		e.capture.Close()
	}
}

func newReadCloser(r io.Reader, closeFunc func() error) io.ReadCloser {
	return &readCloserImpl{
		reader:    r,
		closeFunc: closeFunc,
	}
}

type readCloserImpl struct {
	reader    io.Reader
	closeFunc func() error
}

func (r *readCloserImpl) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

func (r *readCloserImpl) Close() error {
	if r.closeFunc != nil {
		return r.closeFunc()
	}
	return nil
}
