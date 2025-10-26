package dbml

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/dbml"
	"github.com/artarts36/db-exporter/internal/shared/iox"
)

type Exporter struct {
	mapper *mapper
}

func NewExporter() *Exporter {
	return &Exporter{
		mapper: &mapper{},
	}
}

func (e *Exporter) ExportPerFile(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	for _, tbl := range params.Schema.Tables.List() {
		err := params.Workspace.Write(ctx, fmt.Sprintf("%s.dbml", tbl.Name.Value), func(buffer iox.Writer) error {
			dbmlFile := &dbml.File{}
			table, refs, err := e.mapper.mapTable(ctx, tbl, params.Schema.Driver)
			if err != nil {
				return fmt.Errorf("failed to map table %q: w", err)
			}
			dbmlFile.Tables = []*dbml.Table{table}
			dbmlFile.Refs = refs

			dbmlFile.Render(buffer)

			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("write table %q to workspace: %w", tbl.Name.Value, err)
		}
	}

	err := params.Workspace.Write(ctx, "enums.dbml", func(buffer iox.Writer) error {
		dbmlFile := &dbml.File{
			Enums: e.mapper.mapEnums(params.Schema),
		}

		dbmlFile.Render(buffer)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("write enums to workspace: %w", err)
	}

	return nil, nil
}

func (e *Exporter) Export(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	err := params.Workspace.Write(ctx, "schema.dbml", func(buffer iox.Writer) error {
		dbmlFile := &dbml.File{
			Tables: make([]*dbml.Table, 0, params.Schema.Tables.Len()),
			Refs:   make([]*dbml.Ref, 0),
			Enums:  e.mapper.mapEnums(params.Schema),
		}

		for _, tbl := range params.Schema.Tables.List() {
			err := e.exportTable(ctx, tbl, params, dbmlFile)
			if err != nil {
				return fmt.Errorf("export table %q: %w", tbl.Name, err)
			}
		}

		return nil
	})

	return nil, err
}

func (e *Exporter) exportTable(
	ctx context.Context,
	tbl *schema.Table,
	params *exporter.ExportParams,
	dbmlFile *dbml.File,
) error {
	table, refs, err := e.mapper.mapTable(ctx, tbl, params.Schema.Driver)
	if err != nil {
		return fmt.Errorf("map table: %w", err)
	}

	dbmlFile.Tables = append(dbmlFile.Tables, table)
	dbmlFile.Refs = append(dbmlFile.Refs, refs...)

	return nil
}
