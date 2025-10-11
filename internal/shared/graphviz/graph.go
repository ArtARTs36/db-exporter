package graphviz

import (
	"context"
	"errors"
	"fmt"
	"github.com/artarts36/specw"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"image"
)

type Graph struct {
	graph    *cgraph.Graph
	graphviz *graphviz.Graphviz
}

func CreateGraph(ctx context.Context) (*Graph, error) {
	gv := graphviz.New()

	gvGraph, err := gv.Graph()
	if err != nil {
		return nil, fmt.Errorf("create graphviz graphv: %w", err)
	}

	return &Graph{
		graph:    gvGraph,
		graphviz: gv,
	}, nil
}

func (g *Graph) SetBackgroundColor(color specw.Color) {
	g.graph.SetBackgroundColor(color.Hex())
}

func (g *Graph) CreateNode(name string) (*Node, error) {
	node, err := g.graph.CreateNode(name)
	if err != nil {
		return nil, fmt.Errorf("create graphviz node %s: %w", name, err)
	}

	node.SafeSet("labeljust", "c", "")

	node.SetShape(cgraph.PlainTextShape)

	return &Node{node: node, graph: g.graph}, nil
}

func (g *Graph) SetFontName(fontName string) error {
	g.graph.SafeSet("fontname", fontName, "")
	return nil
}

func (g *Graph) CreateEdge(edgeName string, startNode *Node, endNode *Node) (*Edge, error) {
	edge, err := g.graph.CreateEdge(edgeName, startNode.node, endNode.node)
	if err != nil {
		return nil, fmt.Errorf("create graphviz edge %s: %w", edgeName, err)
	}

	return &Edge{edge: edge}, nil
}

func (g *Graph) Close() error {
	errs := []error{}

	if err := g.graph.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close graphviz graph: %w", err))
	}

	if err := g.graphviz.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close graphviz: %w", err))
	}

	return errors.Join(errs...)
}

func (g *Graph) WithoutBackground() {
	g.graph.SetBackgroundColor("transparent")
}

func (g *Graph) Build(ctx context.Context) (image.Image, error) {
	img, err := g.graphviz.RenderImage(g.graph)
	if err != nil {
		return nil, fmt.Errorf("render graphviz image: %w", err)
	}

	return img, nil
}
