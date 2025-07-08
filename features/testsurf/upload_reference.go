package testsurf

import (
	"context"
	"errors"
	"io"

	"github.com/vkhobor/go-opencv/image"
)

type UploadReferenceFeature struct {
}

func (i *UploadReferenceFeature) UploadReference(ctx context.Context, file io.Reader) error {
	ref, err := image.ReadImageFromReader(file)
	if err != nil {
		return err
	}

	if ref.Empty() {
		return errors.New("Reference image is empty")
	}
	if cachedReferenceImage != nil {
		cachedReferenceImage.Close()
	}
	cachedReferenceImage = &ref

	return nil
}
