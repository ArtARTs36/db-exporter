package goentity

import (
	"context"
	"errors"
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/shared/golang"
	"strings"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/config"
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
		page, genErr := e.entityGenerator.Generate(e.entityMapper.MapEntity(table), pkg)
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

	goSch := e.entityMapper.MapEntities(params.Schema.Tables.List())

	goModule := buildGoModule(ctx, e.goModFinder, spec.GoModule, params.Directory)
	pkg, err := golang.BuildPackage(spec.Package, goModule)
	if err != nil {
		return nil, fmt.Errorf("failed to build go package: %w", err)
	}

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

func (e *EntitiesExporter) selectPackage(params *config.GoEntitiesExportSpec) string {
	if params.Package != "" {
		return strings.ToLower(params.Package)
	}

	return "entities"
}
