package goentity

import (
	"context"
	"errors"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/shared/golang"
)

type EntitiesExporter struct {
	pager           *common.Pager
	entityMapper    *EntityMapper
	entityGenerator *EntityGenerator
	goModFinder     *golang.ModFinder
}

func NewEntitiesExporter(
	pager *common.Pager,
	entityMapper *EntityMapper,
	entityGenerator *EntityGenerator,
	goModFinder *golang.ModFinder,
) *EntitiesExporter {
	return &EntitiesExporter{
		pager:           pager,
		entityMapper:    entityMapper,
		entityGenerator: entityGenerator,
		goModFinder:     goModFinder,
	}
}

func (e *EntitiesExporter) ExportPerFile(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.GoEntitiesExportSpec)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	goModule := buildGoModule(ctx, e.goModFinder, spec.GoModule, params.Directory)

	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())
	pkg, err := golang.BuildPackage(spec.Package, goModule)
	if err != nil {
		return nil, fmt.Errorf("failed to build go package: %w", err)
	}

	for _, table := range params.Schema.Tables.List() {
		page, genErr := e.entityGenerator.GenerateEntity(&GenerateEntityParams{
			Entity:  e.entityMapper.MapEntity(table, pkg),
			Package: pkg,
		})
		if genErr != nil {
			return nil, genErr
		}

		pages = append(pages, page)
	}

	return pages, nil
}

func (e *EntitiesExporter) Export(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.GoEntitiesExportSpec)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	goModule := buildGoModule(ctx, e.goModFinder, spec.GoModule, params.Directory)
	pkg, err := buildEntityPackage(spec.Package, goModule)
	if err != nil {
		return nil, fmt.Errorf("failed to build go package: %w", err)
	}

	entities := e.entityMapper.MapEntities(params.Schema.Tables.List(), pkg)

	page, err := e.entityGenerator.GenerateEntities(entities, pkg)
	if err != nil {
		return nil, err
	}

	return []*exporter.ExportedPage{
		page,
	}, nil
}
