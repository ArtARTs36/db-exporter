package graphviz

import (
	"github.com/goccy/go-graphviz/cgraph"
)

type Edge struct {
	edge *cgraph.Edge
}

func (e *Edge) WriteText(txt string) {
	e.edge.SetLabel(txt)
}

func (e *Edge) SetFontSize(size float64) {
	e.edge.SetFontSize(size)
}

func (e *Edge) SetFontName(fontName string) error {
	e.edge.SafeSet("fontname", fontName, "")
	return nil
}
