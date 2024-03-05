package exporter

import (
	"context"
	"fmt"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
	"github.com/huandu/go-sqlbuilder"
	"github.com/mjibson/sqlfmt"
	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/goose"
	"github.com/artarts36/db-exporter/internal/template"
)

type GooseExporter struct {
	renderer *template.Renderer
}

type gooseMigration struct {
	upQueries   []string
	downQueries []string
}

func NewGooseExporter(renderer *template.Renderer) *GooseExporter {
	return &GooseExporter{
		renderer: renderer,
	}
}

func (e *GooseExporter) ExportPerFile(_ context.Context, sch *schema.Schema, params *ExportParams) ([]*ExportedPage, error) {
	pages := make([]*ExportedPage, 0, len(sch.Tables))

	for _, table := range sch.Tables {
		if params.WithoutMigrationsTable && goose.IsMigrationsTable(table.Name.Value) {
			continue
		}

		migration, err := e.makeMigration(table)
		if err != nil {
			return nil, fmt.Errorf("making migration queries failed: %w", err)
		}

		p, err := render(
			e.renderer,
			"goose/migration.sql",
			goose.CreateMigrationFilename(fmt.Sprintf(
				"create_%s_table",
				table.Name.Value,
			)),
			map[string]stick.Value{
				"up_queries":   migration.upQueries,
				"down_queries": migration.downQueries,
			},
		)
		if err != nil {
			return nil, err
		}

		pages = append(pages, p)
	}

	return pages, nil
}

func (e *GooseExporter) Export(_ context.Context, sch *schema.Schema, params *ExportParams) ([]*ExportedPage, error) {
	upQueries := make([]string, 0, len(sch.Tables))
	downQueries := make([]string, 0, len(sch.Tables))

	for _, table := range sch.Tables {
		if params.WithoutMigrationsTable && goose.IsMigrationsTable(table.Name.Value) {
			continue
		}

		migration, err := e.makeMigration(table)
		if err != nil {
			return nil, fmt.Errorf("making migration queries failed: %w", err)
		}

		upQueries = append(upQueries, migration.upQueries...)
		downQueries = append(downQueries, migration.downQueries...)
	}

	p, err := render(
		e.renderer,
		"goose/migration.sql",
		goose.CreateMigrationFilename("init"),
		map[string]stick.Value{
			"up_queries":   upQueries,
			"down_queries": downQueries,
		},
	)
	if err != nil {
		return nil, err
	}

	return []*ExportedPage{
		p,
	}, nil
}

func (e *GooseExporter) makeMigration(table *schema.Table) (*gooseMigration, error) {
	upQuery := sqlbuilder.CreateTable(table.Name.Value)

	for _, column := range table.Columns {
		defs := []string{
			column.Name.Value,
			column.Type.Value,
		}

		upQuery.Define(defs...)
	}

	upSql, err := sqlfmt.FmtSQL(tree.PrettyCfg{
		LineWidth:                4,
		TabWidth:                 4,
		DoNotNewLineAfterColName: true,
	}, []string{upQuery.String()})
	if err != nil {
		return nil, fmt.Errorf("failed to format up query: %w", err)
	}

	return &gooseMigration{
		upQueries: []string{
			upSql,
		},
		downQueries: []string{
			fmt.Sprintf("DROP TABLE %s;", table.Name.Value),
		},
	}, nil
}
