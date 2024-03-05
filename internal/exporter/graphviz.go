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

func buildGraphviz(renderer *template.Renderer, tables map[schema.String]*schema.Table) ([]byte, error) {
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

	log.Print("[graphviz] mapping graph")

	tablesNodes := map[string]*cgraph.Node{}
	for _, table := range tables {
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

	for _, table := range tables {
		tableNode, tnExists := tablesNodes[table.Name.Value]
		if !tnExists {
			continue
		}

		for _, col := range table.Columns {
			if !col.HasForeignKey() {
				continue
			}

			foreignTableNode, ftnExists := tablesNodes[col.ForeignKey.Table.Value]
			if !ftnExists {
				continue
			}

			edge, err := graph.CreateEdge(col.ForeignKey.Name.Value, tableNode, foreignTableNode)
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

	log.Println("[graphviz] generating svg diagram")

	var buf bytes.Buffer
	if err := g.Render(graph, "svg", &buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
