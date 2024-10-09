package exporter

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/template"
)

type GoEntitiesExporter struct {
	renderer     *template.Renderer
	entityMapper *GoEntityMapper
}

func NewGoEntitiesExporter(renderer *template.Renderer, entityMapper *GoEntityMapper) Exporter {
	return &GoEntitiesExporter{
		renderer:     renderer,
		entityMapper: entityMapper,
	}
}

func (e *GoEntitiesExporter) ExportPerFile(
	_ context.Context,
	params *ExportParams,
) ([]*ExportedPage, error) {
	spec, ok := params.Spec.(*config.GoEntitiesExportSpec)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	pages := make([]*ExportedPage, 0, params.Schema.Tables.Len())
	pkg := e.selectPackage(spec)

	for _, table := range params.Schema.Tables.List() {
		goSch := e.entityMapper.MapEntities([]*schema.Table{
			table,
		})

		page, err := render(
			e.renderer,
			"go-entities/model.go.tpl",
			fmt.Sprintf("%s.go", table.Name.Singular().Lower()),
			map[string]stick.Value{
				"schema":  goSch,
				"package": pkg,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to render: %w", err)
		}

		pages = append(pages, page)
	}

	return pages, nil
}

func (e *GoEntitiesExporter) Export(
	_ context.Context,
	params *ExportParams,
) ([]*ExportedPage, error) {
	spec, ok := params.Spec.(*config.GoEntitiesExportSpec)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	goSch := e.entityMapper.MapEntities(params.Schema.Tables.List())
	pkg := e.selectPackage(spec)

	page, err := render(e.renderer, "go-entities/entity.go.tpl", "entities.go", map[string]stick.Value{
		"schema":  goSch,
		"package": pkg,
	})
	if err != nil {
		return nil, err
	}

	return []*ExportedPage{
		page,
	}, nil
}

func (e *GoEntitiesExporter) selectPackage(params *config.GoEntitiesExportSpec) string {
	if params.Package != "" {
		return strings.ToLower(params.Package)
	}

	return "entities"
}
