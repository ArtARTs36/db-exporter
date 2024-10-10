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
	"strings"
)

type RepositoryExporter struct {
	pager           *common.Pager
	goModFinder     *golang.ModFinder
	entityMapper    *EntityMapper
	entityGenerator *EntityGenerator
}

func NewRepositoryExporter(
	pager *common.Pager,
	goModFinder *golang.ModFinder,
	entityMapper *EntityMapper,
	entityGenerator *EntityGenerator,
) *RepositoryExporter {
	return &RepositoryExporter{
		pager:           pager,
		goModFinder:     goModFinder,
		entityMapper:    entityMapper,
		entityGenerator: entityGenerator,
	}
}

type Repository struct {
	Name       string
	Entity     *Entity
	EntityCall string
}

func (e *RepositoryExporter) ExportPerFile(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	const pageTypes = 2

	spec, ok := params.Spec.(*config.GoEntityRepositorySpec)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	goModule := buildGoModule(ctx, e.goModFinder, spec.GoModule, params.Directory)

	pkg, err := e.buildRepositoryPackage(spec, goModule)
	if err != nil {
		return nil, fmt.Errorf("failed to build repository package: %w", err)
	}

	entityPkg, err := buildEntityPackage(spec.Entities.Package, goModule)
	if err != nil {
		return nil, fmt.Errorf("failed to build entity package: %w", err)
	}

	pagesLen := params.Schema.Tables.Len() * pageTypes
	if spec.Repositories.Container.StructName != "" {
		pagesLen++
	}

	pages := make([]*exporter.ExportedPage, 0, pagesLen)
	repositories := make([]*Repository, 0, params.Schema.Tables.Len())

	repoPage := e.pager.Of("go-entities/repository.go.tpl")
	repoNameMaxLength := 0

	for _, table := range params.Schema.Tables.List() {
		entity := e.entityMapper.MapEntity(table)

		repository := &Repository{
			Name:       fmt.Sprintf("PG%sRepository", entity.Name),
			Entity:     entity,
			EntityCall: entityPkg.CallToStruct(pkg, entity.Name.Value),
		}

		if len(repository.Name) > repoNameMaxLength {
			repoNameMaxLength = len(repository.Name)
		}

		repositories = append(repositories, repository)

		page, eerr := e.entityGenerator.GenerateEntity(entity, entityPkg)
		if eerr != nil {
			return nil, fmt.Errorf("failed to generate entity %q: %w", entity.Name, eerr)
		}

		pages = append(pages, page)

		page, rerr := repoPage.Export(
			fmt.Sprintf("%s/%s.go", pkg.ProjectRelativePath, table.Name.Singular().Lower().Value),
			map[string]stick.Value{
				"entityPackage": entityPkg,
				"package":       pkg,
				"schema": map[string]interface{}{
					"Repositories":      []*Repository{repository},
					"RepoNameMaxLength": repoNameMaxLength,
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
				"package":       pkg,
				"containerName": ds.NewString(spec.Repositories.Container.StructName).Pascal().String(),
				"schema": map[string]interface{}{
					"Repositories":      repositories,
					"RepoNameMaxLength": repoNameMaxLength,
				},
			})
		if rerr != nil {
			return nil, rerr
		}
		pages = append(pages, page)
	}

	return pages, nil
}

func (e *RepositoryExporter) Export(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	return e.ExportPerFile(ctx, params)
}

func (e *RepositoryExporter) buildRepositoryPackage(
	spec *config.GoEntityRepositorySpec,
	goModule string,
) (golang.Package, error) {
	pkgName := "repositories"
	if spec.Repositories.Package != "" {
		pkgName = spec.Repositories.Package
	}

	return golang.BuildPackage(pkgName, goModule)
}
