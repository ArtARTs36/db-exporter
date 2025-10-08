package imageencoder

import (
	"image"
	"io"
)

type Encoder interface {
	Encode(writer io.Writer, img image.Image, opts Options) error
}

type Options struct {
	CompressionLevel CompressionLevel
}

type CompressionLevel string

const (
	CompressionLevelUnspecified CompressionLevel = ""
	CompressionLevelNone        CompressionLevel = "none"
	CompressionLevelLow         CompressionLevel = "low"
	CompressionLevelMedium      CompressionLevel = "medium"
	CompressionLevelHigh        CompressionLevel = "high"
)
