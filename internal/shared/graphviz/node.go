package graphviz

import (
	"github.com/goccy/go-graphviz/cgraph"
)

type Node struct {
	node  *cgraph.Node
	graph *cgraph.Graph
}

func (g *Node) SetFontSize(size float64) {
	g.node.SetFontSize(size)
}

func (g *Node) SetFontName(fontName string) error {
	g.node.SafeSet("fontname", fontName, "")
	return nil
}

func (g *Node) WriteHTML(plainHTML string) error {
	preparedHTML := g.graph.StrdupHTML(plainHTML)

	g.node.SetLabel(preparedHTML)

	return nil
}
