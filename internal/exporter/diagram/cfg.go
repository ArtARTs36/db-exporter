package diagram

import (
	"fmt"

	"github.com/artarts36/specw"
	"golang.org/x/image/colornames"
)

type ImageFormat string

const (
	FormatUnspecified ImageFormat = ""
	FormatSVG         ImageFormat = "svg"
	ImageFormatPNG    ImageFormat = "png"
)

type Specification struct {
	Image struct {
		Format ImageFormat `yaml:"format" json:"format"`
	} `yaml:"image" json:"image"`
	Style struct {
		Background struct {
			Grid *struct {
				LineColor *specw.Color `yaml:"line_color" json:"line_color"`
				CellSize  int          `yaml:"cell_size" json:"cell_size"`
			} `yaml:"grid" json:"grid"`
			Color *specw.Color `yaml:"color" json:"color"`
		} `yaml:"background" json:"background"`
		Table struct {
			Name struct {
				BackgroundColor string `yaml:"background_color" json:"background_color"` // #hex
				TextColor       string `yaml:"text_color" json:"text_color"`             // #hex
			} `yaml:"name" json:"name"`
		} `yaml:"table" json:"table"`
		Font struct {
			Family string  `yaml:"family" json:"family"`
			Size   float64 `yaml:"size" json:"size"`
		} `yaml:"font" json:"font"`
	} `yaml:"style" json:"style"`
}

func (s *Specification) Validate() error {
	const (
		defaultGridCellSize = 30
		defaultFontSize     = 32
	)

	if s.Image.Format == FormatUnspecified {
		s.Image.Format = FormatSVG
	} else if !s.Image.Format.Valid() {
		return fmt.Errorf("unknown image format: %s", s.Image.Format)
	}

	if s.Style.Table.Name.BackgroundColor == "" {
		s.Style.Table.Name.BackgroundColor = "#3498db"
	}

	if s.Style.Table.Name.TextColor == "" {
		s.Style.Table.Name.TextColor = "white"
	}

	if s.Style.Background.Color == nil {
		s.Style.Background.Color = &specw.Color{
			Color: colornames.White,
		}
	}

	if s.Style.Background.Grid != nil {
		if s.Style.Background.Grid.LineColor == nil {
			defColor := &specw.Color{}
			defColor.AsEEE()
			s.Style.Background.Grid.LineColor = defColor
		}
		if s.Style.Background.Grid.CellSize == 0 {
			s.Style.Background.Grid.CellSize = defaultGridCellSize
		}
	}

	if s.Style.Font.Size == 0 {
		s.Style.Font.Size = defaultFontSize
	}

	if s.Style.Font.Family == "" {
		s.Style.Font.Family = "Arial"
	}

	return nil
}

func (f ImageFormat) Valid() bool {
	return f == FormatSVG || f == ImageFormatPNG
}
