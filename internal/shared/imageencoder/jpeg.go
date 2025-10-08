package imageencoder

import (
	"image"
	"image/jpeg"
	"io"
)

type JPEG struct{}

func (j *JPEG) Encode(w io.Writer, img image.Image, opts Options) error {
	return jpeg.Encode(w, img, j.mapOptions(opts))
}

func (j *JPEG) mapOptions(opts Options) *jpeg.Options {
	return &jpeg.Options{
		Quality: j.mapQuality(opts.CompressionLevel),
	}
}

func (*JPEG) mapQuality(lvl CompressionLevel) int {
	switch lvl {
	case CompressionLevelNone, CompressionLevelUnspecified:
		return 100 //nolint:mnd //not need
	case CompressionLevelLow:
		return 70 //nolint:mnd //not need
	case CompressionLevelMedium:
		return 40 //nolint:mnd //not need
	case CompressionLevelHigh:
		return 10 //nolint:mnd //not need
	}
	return 100 //nolint:mnd //not need
}
