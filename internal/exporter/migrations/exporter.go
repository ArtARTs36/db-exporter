package migrations

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/infrastructure/sql"
)

type Exporter struct {
	name string

	page              *common.Page
	maker             MigrationMaker
	ddlBuilderManager *sql.DDLBuilderManager
}

func NewExporter(
	name string,
	pager *common.Pager,
	templateName string,
	ddlBuilder *sql.DDLBuilderManager,
	maker MigrationMaker,
) *Exporter {
	return &Exporter{
		name:              name,
		page:              pager.Of(templateName),
		maker:             maker,
		ddlBuilderManager: ddlBuilder,
	}
}

func (e *Exporter) ExportPerFile(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*Specification)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())

	slog.DebugContext(ctx, fmt.Sprintf("[%s] building queries and rendering migration files", e.name))

	ddlOpts := e.mapBuildDDLOpts(spec)

	ddlBuilder := e.ddlBuilderManager.For(spec.Target)
	ddls, err := ddlBuilder.BuildPerTable(params.Schema, ddlOpts)
	if err != nil {
		return nil, err
	}

	for i, ddl := range ddls {
		migMeta := e.maker.MakeSingle(i, ddl.Name)
		mig := &Migration{
			Meta:        migMeta.Attrs,
			UpQueries:   ddl.UpQueries,
			DownQueries: ddl.DownQueries,
		}

		p, expErr := e.page.Export(migMeta.Filename, map[string]stick.Value{
			"migration": mig,
		})
		if expErr != nil {
			return nil, expErr
		}

		pages = append(pages, p)
	}

	return pages, nil
}

func (e *Exporter) mapBuildDDLOpts(spec *Specification) sql.BuildDDLOpts {
	return sql.BuildDDLOpts{
		UseIfNotExists: spec.Use.IfNotExists,
		UseIfExists:    spec.Use.IfExists,
	}
}

func (e *Exporter) Export(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*Specification)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	slog.DebugContext(ctx, fmt.Sprintf("[%s] building queries", e.name))

	ddlBuilder := e.ddlBuilderManager.For(spec.Target)
	ddl, err := ddlBuilder.Build(params.Schema, e.mapBuildDDLOpts(spec))
	if err != nil {
		return nil, fmt.Errorf("failed to build ddl for schema: %w", err)
	}

	migMeta := e.maker.MakeMultiple()
	mig := &Migration{
		Meta:        migMeta.Attrs,
		UpQueries:   ddl.UpQueries,
		DownQueries: ddl.DownQueries,
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
