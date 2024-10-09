package goentity

import (
	"context"
	"errors"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/shared/ds"
	"github.com/artarts36/db-exporter/internal/shared/golang"
	"github.com/tyler-sommer/stick"
	"log/slog"
	"strings"
)

type GoEntityRepositoryExporter struct {
	pager              *common.Pager
	goModFinder        *golang.ModFinder
	entityMapper       *EntityMapper
	goEntitiesExporter *EntitiesExporter
}

func NewRepositoryExporter(
	pager *common.Pager,
	goModFinder *golang.ModFinder,
	entityMapper *EntityMapper,
	exporter *EntitiesExporter,
) *GoEntityRepositoryExporter {
	return &GoEntityRepositoryExporter{
		pager:              pager,
		goModFinder:        goModFinder,
		entityMapper:       entityMapper,
		goEntitiesExporter: exporter,
	}
}

type GoEntityRepository struct {
	Name       string
	Entity     *Entity
	EntityCall string
}

func (e *GoEntityRepositoryExporter) ExportPerFile(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.GoEntityRepositorySpec)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	goModule := spec.GoModule
	if goModule == "" {
		goMod, err := e.goModFinder.FindIn(params.Directory)
		if err != nil {
			slog.
				With(slog.Any("err", err)).
				WarnContext(ctx, "[go-entity-repository-exporter] failed to get go module")
		} else {
			goModule = goMod.Module
		}
	}

	pkg, err := e.buildRepositoryPackage(spec, goModule)
	if err != nil {
		return nil, fmt.Errorf("failed to build repository package: %w", err)
	}

	entityPkg, err := e.buildEntityPackage(spec, goModule)
	if err != nil {
		return nil, fmt.Errorf("failed to build entity package: %w", err)
	}

	pagesLen := params.Schema.Tables.Len() * 2
	if spec.Repositories.Container.StructName != "" {
		pagesLen++
	}

	pages := make([]*exporter.ExportedPage, 0, pagesLen)
	repositories := make([]*GoEntityRepository, 0, params.Schema.Tables.Len())

	repoPage := e.pager.Of("go-entities/repository.go.tpl")

	for _, table := range params.Schema.Tables.List() {
		entity := e.entityMapper.MapEntity(table)

		repository := &GoEntityRepository{
			Name:       fmt.Sprintf("PG%sRepository", entity.Name),
			Entity:     entity,
			EntityCall: entityPkg.CallToStruct(pkg, entity.Name.Value),
		}
		repositories = append(repositories, repository)

		page, eerr := e.goEntitiesExporter.generateEntity(entity, entityPkg.Name)
		if eerr != nil {
			return nil, fmt.Errorf("failed to generate entity %q: %w", entity.Name, eerr)
		}

		pages = append(pages, page)

		page, rerr := repoPage.Export(
			fmt.Sprintf("%s/%s.go", pkg.ProjectRelativePath, table.Name.Singular().Lower().Value),
			map[string]stick.Value{
				"package": pkg,
				"schema": map[string]interface{}{
					"Repositories": []*GoEntityRepository{repository},
				},
			},
		)
		if rerr != nil {
			return nil, rerr
		}
		pages = append(pages, page)
	}

	if spec.Repositories.Container.StructName != "" {
		page, rerr := e.pager.Of("go-entities/container.go.tpl").Export(
			fmt.Sprintf("%s/%s.go", pkg.ProjectRelativePath, strings.ToLower(spec.Repositories.Container.StructName)),
			map[string]stick.Value{
				"package":        pkg,
				"container_name": ds.NewString(spec.Repositories.Container.StructName).Pascal().String(),
				"schema": map[string]interface{}{
					"Repositories": repositories,
				},
			})
		if rerr != nil {
			return nil, rerr
		}
		pages = append(pages, page)
	}

	return pages, nil
}

func (e *GoEntityRepositoryExporter) Export(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	return e.ExportPerFile(ctx, params)
}

func (e *GoEntityRepositoryExporter) buildRepositoryPackage(
	spec *config.GoEntityRepositorySpec,
	goModule string,
) (golang.Package, error) {
	pkgName := "repositories"
	if spec.Repositories.Package != "" {
		pkgName = spec.Repositories.Package
	}

	return golang.BuildPackage(pkgName, goModule)
}

func (e *GoEntityRepositoryExporter) buildEntityPackage(
	spec *config.GoEntityRepositorySpec,
	goModule string,
) (golang.Package, error) {
	pkgName := "entities"
	if spec.Entities.Package != "" {
		pkgName = spec.Entities.Package
	}

	return golang.BuildPackage(pkgName, goModule)
}
