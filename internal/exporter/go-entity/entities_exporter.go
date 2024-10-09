package goentity

import (
	"context"
	"errors"
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"strings"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/config"
)

type EntitiesExporter struct {
	pager        *common.Pager
	entityMapper *EntityMapper
}

func NewEntitiesExporter(pager *common.Pager, entityMapper *EntityMapper) *EntitiesExporter {
	return &EntitiesExporter{
		pager:        pager,
		entityMapper: entityMapper,
	}
}

func (e *EntitiesExporter) ExportPerFile(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.GoEntitiesExportSpec)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())
	pkg := e.selectPackage(spec)

	for _, table := range params.Schema.Tables.List() {
		page, err := e.generateEntity(e.entityMapper.MapEntity(table), pkg)
		if err != nil {
			return nil, err
		}

		pages = append(pages, page)
	}

	return pages, nil
}

func (e *EntitiesExporter) Export(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.GoEntitiesExportSpec)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	goSch := e.entityMapper.MapEntities(params.Schema.Tables.List())
	pkg := e.selectPackage(spec)

	page, err := e.pager.Of("go-entities/entity.go.tpl").Export("entities.go", map[string]stick.Value{
		"schema":  goSch,
		"package": pkg,
	})
	if err != nil {
		return nil, err
	}

	return []*exporter.ExportedPage{
		page,
	}, nil
}

func (e *EntitiesExporter) generateEntity(entity *Entity, pkg string) (*exporter.ExportedPage, error) {
	page, err := e.pager.Of("go-entities/entity.go.tpl").Export(
		fmt.Sprintf("%s.go", entity.Table.Name.Singular().Lower()),
		map[string]stick.Value{
			"schema": &Entities{
				Entities: []*Entity{entity},
				Imports:  entity.Imports,
			},
			"package": pkg,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to render: %w", err)
	}
	return page, nil
}

func (e *EntitiesExporter) selectPackage(params *config.GoEntitiesExportSpec) string {
	if params.Package != "" {
		return strings.ToLower(params.Package)
	}

	return "entities"
}
