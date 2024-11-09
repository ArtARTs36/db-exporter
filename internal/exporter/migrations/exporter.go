package migrations

import (
	"context"
	"errors"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"log/slog"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/sql"
)

type Exporter struct {
	name string

	page       *common.Page
	maker      MigrationMaker
	ddlBuilder *sql.DDLBuilder
}

func NewExporter(
	name string,
	pager *common.Pager,
	templateName string,
	ddlBuilder *sql.DDLBuilder,
	maker MigrationMaker,
) *Exporter {
	return &Exporter{
		name:       name,
		page:       pager.Of(templateName),
		maker:      maker,
		ddlBuilder: ddlBuilder,
	}
}

func (e *Exporter) ExportPerFile(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.MigrationsSpec)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())

	slog.DebugContext(ctx, fmt.Sprintf("[%s] building queries and rendering migration files", e.name))

	ddlOpts := sql.BuildDDLParams{
		UseIfNotExists: spec.Use.IfNotExists,
		Source:         params.Schema.Driver,
		Target:         spec.Target,
	}

	for i, table := range params.Schema.Tables.List() {
		upQueries, downQueries := make([]string, 0), make([]string, 0)

		for _, sequence := range table.UsingSequences {
			if sequence.Used == 1 {
				seqParams := sql.CreateSequenceParams{
					UseIfNotExists: spec.Use.IfExists,
					Source:         params.Schema.Driver,
					Target:         spec.Target,
				}

				seqSQL, err := e.ddlBuilder.CreateSequence(sequence, seqParams)
				if err != nil {
					return nil, fmt.Errorf("failed to build query for create sequence %q: %w", sequence.Name, err)
				}

				upQueries = append(upQueries, seqSQL)
				downQueries = append(downQueries, e.ddlBuilder.DropSequence(sequence, spec.Use.IfExists))
			}
		}

		for _, enum := range table.UsingEnums {
			if enum.Used == 1 {
				upQueries = append(upQueries, e.ddlBuilder.CreateEnum(enum))
				downQueries = append(downQueries, e.ddlBuilder.DropType(enum.Name.Value, spec.Use.IfExists))
			}
		}

		upQ, downQ, err := e.createQueries(table, ddlOpts, spec.Use.IfExists)
		if err != nil {
			return nil, fmt.Errorf("failed to create queries: %w", err)
		}

		upQueries = append(upQueries, upQ...)
		downQueries = append(downQueries, downQ...)

		migMeta := e.maker.MakeSingle(i, table.Name)
		mig := &Migration{
			Meta:        migMeta.Attrs,
			UpQueries:   upQueries,
			DownQueries: downQueries,
		}

		p, err := e.page.Export(migMeta.Filename, map[string]stick.Value{
			"migration": mig,
		})
		if err != nil {
			return nil, err
		}

		pages = append(pages, p)
	}

	return pages, nil
}

func (e *Exporter) Export(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.MigrationsSpec)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	upQueries := make([]string, 0, params.Schema.Tables.Len())
	downQueries := make([]string, 0, params.Schema.Tables.Len())

	slog.DebugContext(ctx, fmt.Sprintf("[%s] building queries", e.name))

	ddlOpts := sql.BuildDDLParams{
		UseIfNotExists: spec.Use.IfNotExists,
		Source:         params.Schema.Driver,
		Target:         spec.Target,
	}

	for _, table := range params.Schema.Tables.List() {
		upQs, downQs, err := e.createQueries(table, ddlOpts, spec.Use.IfExists)
		if err != nil {
			return nil, fmt.Errorf("failed to create queries: %w", err)
		}

		upQueries = append(upQueries, upQs...)
		downQueries = append(downQueries, downQs...)
	}

	seqParams := sql.CreateSequenceParams{
		UseIfNotExists: spec.Use.IfExists,
		Source:         params.Schema.Driver,
		Target:         spec.Target,
	}

	for _, seq := range params.Schema.Sequences {
		seqSQL, err := e.ddlBuilder.CreateSequence(seq, seqParams)
		if err != nil {
			return nil, fmt.Errorf("failed to build query for create sequence %q: %w", seq.Name, err)
		}

		upQueries = append(upQueries, seqSQL)
		downQueries = append(downQueries, e.ddlBuilder.DropSequence(seq, spec.Use.IfNotExists))
	}

	migMeta := e.maker.MakeMultiple()
	mig := &Migration{
		Meta:        migMeta.Attrs,
		UpQueries:   upQueries,
		DownQueries: downQueries,
	}

	p, err := e.page.Export(migMeta.Filename, map[string]stick.Value{
		"migration": mig,
	})
	if err != nil {
		return nil, err
	}

	return []*exporter.ExportedPage{
		p,
	}, nil
}

func (e *Exporter) createQueries(
	table *schema.Table,
	opts sql.BuildDDLParams,
	useIfExists bool,
) (
	upQueries []string,
	downQueries []string,
	err error,
) {
	upQueries, err = e.ddlBuilder.BuildDDL(table, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build ddl: %w", err)
	}

	downQueries = []string{
		e.ddlBuilder.DropTable(table, useIfExists),
	}
	return //nolint: nakedret // not need
}
