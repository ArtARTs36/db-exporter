package graphviz

import (
	"fmt"
	"github.com/goccy/go-graphviz"
)

type Node struct {
	node  *graphviz.Node
	graph *graphviz.Graph
}

func (g *Node) SetFontName(fontName string) error {
	return g.node.SafeSet("fontname", fontName, "")
}

func (g *Node) WriteHTML(plainHTML string) error {
	preparedHTML, err := g.graph.StrdupHTML(plainHTML)
	if err != nil {
		return fmt.Errorf("prepare HTML: %w", err)
	}

	g.node.SetLabel(preparedHTML)

	return nil
}
