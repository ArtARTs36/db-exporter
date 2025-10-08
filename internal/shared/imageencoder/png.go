package imageencoder

import (
	"image"
	"image/png"
	"io"
)

type PNG struct{}

func (p *PNG) Encode(w io.Writer, img image.Image, opts Options) error {
	encoder := &png.Encoder{
		CompressionLevel: p.mapCompressionLevel(opts.CompressionLevel),
	}

	return encoder.Encode(w, img)
}

func (*PNG) mapCompressionLevel(lvl CompressionLevel) png.CompressionLevel {
	switch lvl {
	case CompressionLevelNone, CompressionLevelUnspecified:
		return png.NoCompression
	case CompressionLevelLow:
		return png.BestSpeed
	case CompressionLevelMedium:
		return png.DefaultCompression
	case CompressionLevelHigh:
		return png.BestCompression
	}
	return png.NoCompression
}
