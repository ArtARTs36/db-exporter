package imageencoder

type Format string

const (
	FormatUnspecified Format = ""
	FormatPNG         Format = "png"
	FormatJPEG        Format = "jpeg"
)

func (f Format) Valid() bool {
	return f == FormatUnspecified || f == FormatPNG || f == FormatJPEG
}
