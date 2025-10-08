package diagram

import (
	"bytes"
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/schema"
	graphviz2 "github.com/artarts36/db-exporter/internal/shared/graphviz"
	"github.com/artarts36/db-exporter/internal/template"
	"github.com/tyler-sommer/stick"
	"log/slog"
)

type GraphBuilder struct {
	renderer *template.Renderer
}

func NewGraphBuilder(renderer *template.Renderer) *GraphBuilder {
	return &GraphBuilder{renderer: renderer}
}

func (b *GraphBuilder) BuildSVG(tables *schema.TableMap, spec *config.DiagramExportSpec) ([]byte, error) {
	graph, err := b.buildGraph(tables, spec)
	if err != nil {
		return nil, fmt.Errorf("build graph: %w", err)
	}

	defer func() {
		if err = graph.Close(); err != nil {
			slog.Warn("failed to close graph", slog.String("err", err.Error()))
		}
	}()

	slog.Debug("[diagram] generating svg diagram")

	var buf bytes.Buffer
	if err = graph.Render(context.Background(), "svg", &buf); err != nil {
		return nil, fmt.Errorf("to render grapgh to svg: %w", err)
	}

	return buf.Bytes(), nil
}

func (b *GraphBuilder) buildGraph(
	tables *schema.TableMap,
	spec *config.DiagramExportSpec,
) (*graphviz2.Graph, error) {
	graph, err := graphviz2.CreateGraph(context.Background())
	if err != nil {
		return graph, fmt.Errorf("failed to create graph: %w", err)
	}

	graph.SetBackgroundColor(spec.Style.Background.Color)

	slog.Debug("[graphbuilder] mapping graph")

	tablesNodes, err := b.buildNodes(graph, tables, spec)
	if err != nil {
		return nil, fmt.Errorf("failed to build nodes: %w", err)
	}

	slog.Debug(fmt.Sprintf("[graphbuilder] builded %d nodes", len(tablesNodes)))

	edgesCount, err := b.buildEdges(graph, tables, tablesNodes, spec)
	if err != nil {
		return nil, fmt.Errorf("failed to build edges: %w", err)
	}

	slog.Debug(fmt.Sprintf("[graphbuilder] builded %d edges", edgesCount))

	return graph, nil
}

func (b *GraphBuilder) buildNodes(
	graph *graphviz2.Graph,
	tables *schema.TableMap,
	spec *config.DiagramExportSpec,
) (map[string]*graphviz2.Node, error) {
	tablesNodes := map[string]*graphviz2.Node{}

	err := tables.EachWithErr(func(table *schema.Table) error {
		node, graphErr := graph.CreateNode(table.Name.Value)
		if graphErr != nil {
			return fmt.Errorf("failed to create node for table %q: %w", table.Name.Value, graphErr)
		}

		if err := node.SetFontName(spec.Style.Font.Family); err != nil {
			return fmt.Errorf("set font name: %w", err)
		}

		ht, tableErr := b.renderer.Render("@embed/diagram/table.html", map[string]stick.Value{
			"table": mapTable(table),
			"style": spec.Style,
		})
		if tableErr != nil {
			return tableErr
		}

		if err := node.WriteHTML(string(ht)); err != nil {
			return fmt.Errorf("write html to node: %w", err)
		}

		tablesNodes[table.Name.Value] = node

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tablesNodes, nil
}

func (b *GraphBuilder) buildEdges(
	graph *graphviz2.Graph,
	tables *schema.TableMap,
	tablesNodes map[string]*graphviz2.Node,
	spec *config.DiagramExportSpec,
) (int, error) {
	edges := 0

	err := tables.EachWithErr(func(table *schema.Table) error {
		tableNode, tnExists := tablesNodes[table.Name.Value]
		if !tnExists {
			return nil
		}

		for _, col := range table.Columns {
			if !col.HasForeignKey() {
				continue
			}

			foreignTableNode, ftnExists := tablesNodes[col.ForeignKey.ForeignTable.Value]
			if !ftnExists {
				continue
			}

			edge, edgeErr := graph.CreateEdge(col.ForeignKey.Name.Value, tableNode, foreignTableNode)
			if edgeErr != nil {
				return fmt.Errorf(
					"failed to create edge from %s.%s to %s.%s: %w",
					table.Name.Value,
					col.Name.Value,
					col.ForeignKey.ForeignTable,
					col.ForeignKey.ForeignColumn,
					edgeErr,
				)
			}

			edges++

			edge.WriteText(fmt.Sprintf("  %s:%s", col.Name.Value, col.ForeignKey.ForeignColumn.Value))

			if err := edge.SetFontName(spec.Style.Font.Family); err != nil {
				return fmt.Errorf("set font name for edge: %w", err)
			}
		}

		return nil
	})

	return edges, err
}
