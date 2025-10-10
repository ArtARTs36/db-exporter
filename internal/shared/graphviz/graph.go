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
	graph    *graphviz.Graph
	graphviz *graphviz.Graphviz
}

func CreateGraph(ctx context.Context) (*Graph, error) {
	gv, err := graphviz.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("create graphviz: %w", err)
	}

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
	node, err := g.graph.CreateNodeByName(name)
	if err != nil {
		return nil, fmt.Errorf("create graphviz node %s: %w", name, err)
	}

	node.SetShape(cgraph.PlainTextShape)
	if err = node.SafeSet("labeljust", "c", ""); err != nil {
		return nil, fmt.Errorf("set labeljust: %w", err)
	}

	return &Node{node: node, graph: g.graph}, nil
}

func (g *Graph) SetFontName(fontName string) error {
	return g.graph.SafeSet("fontname", fontName, "")
}

func (g *Graph) CreateEdge(edgeName string, startNode *Node, endNode *Node) (*Edge, error) {
	edge, err := g.graph.CreateEdgeByName(edgeName, startNode.node, endNode.node)
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
	img, err := g.graphviz.RenderImage(ctx, g.graph)
	if err != nil {
		return nil, fmt.Errorf("render graphviz image: %w", err)
	}

	return img, nil
}
