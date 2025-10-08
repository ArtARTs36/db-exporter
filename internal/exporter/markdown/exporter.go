package markdown

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/diagram"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/gds"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/schema"
)

type Exporter struct {
	pager          *common.Pager
	diagramCreator *diagram.Creator
}

type markdownPreparedTable struct {
	*schema.Table
	FileName string
}

func NewMarkdownExporter(pager *common.Pager, diagramCreator *diagram.Creator) exporter.Exporter {
	return &Exporter{
		pager:          pager,
		diagramCreator: diagramCreator,
	}
}

func (e *Exporter) ExportPerFile(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.MarkdownExportSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec, expected MarkdownExportSpec, got %T", params.Spec)
	}

	var diag *exporter.ExportedPage
	pagesCap := params.Schema.Tables.Len() + 1
	if spec.WithDiagram {
		pagesCap++
		diagContent, err := e.diagramCreator.Create(ctx, params.Schema.Tables, &spec.Diagram)
		if err != nil {
			return nil, fmt.Errorf("failed to build diag: %w", err)
		}

		diag = &exporter.ExportedPage{
			FileName: fmt.Sprintf("diagram.%s", string(spec.Diagram.Image.Format)),
			Content:  diagContent,
		}
	}

	pages := make([]*exporter.ExportedPage, 0, pagesCap)
	preparedTables := make([]*markdownPreparedTable, 0, params.Schema.Tables.Len())

	mdPage := e.pager.Of("@embed/md/pet-table.md")

	for _, table := range params.Schema.Tables.List() {
		fileName := fmt.Sprintf("%s.md", table.Name.Value)

		page, err := mdPage.Export(fileName, map[string]stick.Value{
			"table": table,
		})
		if err != nil {
			return nil, err
		}

		pages = append(pages, page)

		preparedTables = append(preparedTables, &markdownPreparedTable{
			Table:    table,
			FileName: fileName,
		})
	}

	indexPage, err := e.pager.Of("@embed/md/per-index.md").Export("index.md", map[string]stick.Value{
		"tables":  preparedTables,
		"diagram": diag,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to generate index file: %w", err)
	}

	pages = append(pages, indexPage)

	if diag != nil {
		pages = append(pages, diag)
	}

	return pages, nil
}

func (e *Exporter) Export(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.MarkdownExportSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec, expected MarkdownExportSpec, got %T", params.Spec)
	}

	var diag *exporter.ExportedPage

	if spec.WithDiagram {
		diagContent, err := e.diagramCreator.Create(ctx, params.Schema.Tables, &spec.Diagram)
		if err != nil {
			return nil, fmt.Errorf("build diagram: %w", err)
		}

		diag = &exporter.ExportedPage{
			FileName: fmt.Sprintf("diagram.%s", string(spec.Diagram.Image.Format)),
			Content:  diagContent,
		}
	}

	page, err := e.pager.Of("@embed/md/single-tables.md").Export(
		e.createIndexPageName(params.Schema),
		map[string]stick.Value{
			"schema":        params.Schema,
			"diagram":       diag,
			"diagramExists": diag != nil,
		},
	)
	if err != nil {
		return nil, err
	}

	pages := []*exporter.ExportedPage{
		page,
	}

	if diag != nil {
		pages = append(pages, diag)
	}

	return pages, nil
}

func (e *Exporter) createIndexPageName(sch *schema.Schema) string {
	if sch.Tables.Has(gds.String{Value: "INDEX"}) {
		return "index.md"
	}

	return "INDEX.md"
}
