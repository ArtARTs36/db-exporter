package exporter

import (
	"bytes"
	"fmt"
	"github.com/artarts36/db-exporter/internal/template"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/tyler-sommer/stick"
	"log"

	"github.com/goccy/go-graphviz"

	"github.com/artarts36/db-exporter/internal/schema"
)

func buildGraphviz(renderer *template.Renderer, sc *schema.Schema) ([]byte, error) {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		return nil, fmt.Errorf("failed to create graph: %w", err)
	}

	defer func() {
		if err := graph.Close(); err != nil {
			log.Printf("failed to close graph: %v", err.Error())
		}
		if err := g.Close(); err != nil {
			log.Printf("failed to close graph: %v", err.Error())
		}
	}()

	tablesNodes := map[string]*cgraph.Node{}
	for _, table := range sc.Tables {
		node, err := graph.CreateNode(table.Name.Value)
		if err != nil {
			return nil, err
		}

		node.SetShape(cgraph.PlainTextShape)
		node.SafeSet("class", "db-tables", "")

		ht, err := renderer.Render("graphviz/table.html", map[string]stick.Value{
			"table": table,
		})
		if err != nil {
			return nil, err
		}

		node.SetLabel(graph.StrdupHTML(string(ht)))

		tablesNodes[table.Name.Value] = node
	}

	for _, table := range sc.Tables {
		tableNode, _ := tablesNodes[table.Name.Value]

		for _, col := range table.Columns {
			if !col.HasForeignKey() {
				continue
			}

			edge, err := graph.CreateEdge(col.ForeignKey.Name.Value, tableNode, tablesNodes[col.ForeignKey.Table.Value])
			if err != nil {
				return nil, fmt.Errorf(
					"failed to create edge from %s.%s to %s.%s: %w",
					table.Name.Value,
					col.Name.Value,
					col.ForeignKey.Table,
					col.ForeignKey.Column,
					err,
				)
			}

			edge.SetLabel(fmt.Sprintf(
				"  %s:%s",
				col.Name.Value,
				col.ForeignKey.Column.Value,
			))
		}
	}

	var buf bytes.Buffer
	if err := g.Render(graph, "svg", &buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
