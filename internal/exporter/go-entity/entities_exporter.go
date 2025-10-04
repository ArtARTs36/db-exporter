package goentity

import (
	"context"
	"errors"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/shared/golang"
	"github.com/tyler-sommer/stick"
)

type EntitiesExporter struct {
	pager           *common.Pager
	entityMapper    *EntityMapper
	entityGenerator *EntityGenerator
	goModFinder     *golang.ModFinder

	page struct {
		enumString *common.Page
	}
}

func NewEntitiesExporter(
	pager *common.Pager,
	entityMapper *EntityMapper,
	entityGenerator *EntityGenerator,
	goModFinder *golang.ModFinder,
) *EntitiesExporter {
	exp := &EntitiesExporter{
		pager:           pager,
		entityMapper:    entityMapper,
		entityGenerator: entityGenerator,
		goModFinder:     goModFinder,
	}

	exp.page.enumString = pager.Of("@embed/go-entities/enum_string.go.tpl")

	return exp
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
	pkg, err := buildEntityPackage(spec.Package, goModule)
	if err != nil {
		return nil, fmt.Errorf("failed to build go package: %w", err)
	}

	enums := map[string]*golang.StringEnum{}

	for _, enum := range params.Schema.Enums {
		enumFile := golang.NewFile(fmt.Sprintf("%s.go", enum.Name.Value), pkg)

		goEnum := golang.NewStringEnumOfValues(enum.Name, enum.Values)
		enums[enum.Name.Value] = goEnum

		page, enumErr := e.page.enumString.Export(
			fmt.Sprintf("%s/%s", pkg.ProjectRelativePath, enumFile.Name),
			map[string]stick.Value{
				"enum":  goEnum,
				"_file": enumFile,
			},
		)
		if enumErr != nil {
			return nil, fmt.Errorf("failed to generate enum %q: %w", enum.Name.Value, enumErr)
		}

		pages = append(pages, page)
	}

	for _, table := range params.Schema.Tables.List() {
		page, genErr := e.entityGenerator.GenerateEntity(&GenerateEntityParams{
			Entity: e.entityMapper.MapEntity(&MapEntityParams{
				SourceDriver: params.Schema.Driver,
				Table:        table,
				Package:      pkg,
				Enums:        enums,
			}),
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

	enums := map[string]*golang.StringEnum{}

	for _, enum := range params.Schema.Enums {
		goEnum := golang.NewStringEnumOfValues(enum.Name, enum.Values)
		enums[enum.Name.Value] = goEnum
	}

	entities := e.entityMapper.MapEntities(&MapEntitiesParams{
		SourceDriver: params.Schema.Driver,
		Tables:       params.Schema.Tables.List(),
		Package:      pkg,
		Enums:        map[string]*golang.StringEnum{},
	})

	page, err := e.entityGenerator.GenerateEntities(entities, pkg, enums)
	if err != nil {
		return nil, err
	}

	return []*exporter.ExportedPage{
		page,
	}, nil
}
