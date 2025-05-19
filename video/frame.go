package video

import (
	"bytes"
	"errors"
	"io"
	"os"

	"gocv.io/x/gocv"
)

type FrameExtractor struct {
	capture *gocv.VideoCapture
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
		capture,
	}, nil
}

func (e FrameExtractor) GetFrame(frame int) (io.ReadCloser, error) {
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
	data, err := gocv.IMEncode(".jpg", mat)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	return newReadCloser(bytes.NewReader(data.GetBytes()), func() error {
		data.Close()
		return nil
	}), nil
}

func (e FrameExtractor) Close() {
	e.capture.Close()
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
