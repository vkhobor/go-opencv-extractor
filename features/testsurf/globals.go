package testsurf

import (
	"github.com/vkhobor/go-opencv/video"
	"gocv.io/x/gocv"
)

var cachedTestVideoExtractor video.FrameExtractor
var cachedReferenceImage *gocv.Mat
