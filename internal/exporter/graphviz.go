package exporter

import (
	"bytes"
	"fmt"
	"log"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/template"
)

type graphBuilder struct {
	renderer *template.Renderer
}

func (b *graphBuilder) BuildSVG(tables *schema.TableMap) ([]byte, error) {
	g, graph, err := b.buildGraph(tables)
	if err != nil {
		return nil, fmt.Errorf("failed to build graph: %w", err)
	}

	defer func() {
		if err = graph.Close(); err != nil {
			log.Printf("failed to close graph: %v", err.Error())
		}
		if err = g.Close(); err != nil {
			log.Printf("failed to close graph: %v", err.Error())
		}
	}()

	log.Println("[diagram] generating svg diagram")

	var buf bytes.Buffer
	if err = g.Render(graph, "svg", &buf); err != nil {
		return nil, fmt.Errorf("failed to render grapgh to svg: %w", err)
	}

	return buf.Bytes(), nil
}

func (b *graphBuilder) buildGraph(tables *schema.TableMap) (*graphviz.Graphviz, *cgraph.Graph, error) {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		return g, graph, fmt.Errorf("failed to create graph: %w", err)
	}

	log.Print("[graphbuilder] mapping graph")

	tablesNodes, err := b.buildNodes(graph, tables)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build nodes: %w", err)
	}

	log.Printf("[graphbuilder] builded %d nodes", len(tablesNodes))

	edgesCount, err := b.buildEdges(graph, tables, tablesNodes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build edges: %w", err)
	}

	log.Printf("[graphbuilder] builded %d edges", edgesCount)

	return g, graph, nil
}

func (b *graphBuilder) buildNodes(graph *cgraph.Graph, tables *schema.TableMap) (map[string]*cgraph.Node, error) {
	tablesNodes := map[string]*cgraph.Node{}

	err := tables.EachWithErr(func(table *schema.Table) error {
		node, graphErr := graph.CreateNode(table.Name.Val)
		if graphErr != nil {
			return fmt.Errorf("failed to create node for table %q: %w", table.Name.Val, graphErr)
		}

		node.SetShape(cgraph.PlainTextShape)
		node.SafeSet("class", "db-tables", "")

		ht, tableErr := b.renderer.Render("diagram/table.html", map[string]stick.Value{
			"table": table,
		})
		if tableErr != nil {
			return tableErr
		}

		node.SetLabel(graph.StrdupHTML(string(ht)))

		tablesNodes[table.Name.Val] = node

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tablesNodes, nil
}

func (b *graphBuilder) buildEdges(
	graph *cgraph.Graph,
	tables *schema.TableMap,
	tablesNodes map[string]*cgraph.Node,
) (int, error) {
	edges := 0

	err := tables.EachWithErr(func(table *schema.Table) error {
		tableNode, tnExists := tablesNodes[table.Name.Val]
		if !tnExists {
			return nil
		}

		for _, col := range table.Columns {
			if !col.HasForeignKey() {
				continue
			}

			foreignTableNode, ftnExists := tablesNodes[col.ForeignKey.ForeignTable.Val]
			if !ftnExists {
				continue
			}

			edge, edgeErr := graph.CreateEdge(col.ForeignKey.Name.Val, tableNode, foreignTableNode)
			if edgeErr != nil {
				return fmt.Errorf(
					"failed to create edge from %s.%s to %s.%s: %w",
					table.Name.Val,
					col.Name.Val,
					col.ForeignKey.ForeignTable,
					col.ForeignKey.ForeignColumn,
					edgeErr,
				)
			}

			edges++

			edge.SetLabel(fmt.Sprintf(
				"  %s:%s",
				col.Name.Val,
				col.ForeignKey.ForeignColumn.Val,
			))
		}

		return nil
	})

	return edges, err
}
