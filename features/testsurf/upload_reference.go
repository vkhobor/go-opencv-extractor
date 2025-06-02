package testsurf

import (
	"context"
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

	cachedReferenceImage = ref

	return nil
}
