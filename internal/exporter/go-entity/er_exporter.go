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
	goModFinder     *golang.ModFinder
	entityMapper    *EntityMapper
	entityGenerator *EntityGenerator
	propertyMapper  *GoPropertyMapper

	page struct {
		repo               *common.Page
		entityRepo         *common.Page
		container          *common.Page
		containerInterface *common.Page
	}
}

func NewRepositoryExporter(
	pager *common.Pager,
	goModFinder *golang.ModFinder,
	entityMapper *EntityMapper,
	entityGenerator *EntityGenerator,
	propertyMapper *GoPropertyMapper,
) *RepositoryExporter {
	exp := &RepositoryExporter{
		goModFinder:     goModFinder,
		entityMapper:    entityMapper,
		entityGenerator: entityGenerator,
		propertyMapper:  propertyMapper,
	}

	exp.page.repo = pager.Of("go-entities/repository.go.tpl")
	exp.page.entityRepo = pager.Of("go-entities/entity_repos.go.tpl")
	exp.page.container = pager.Of("go-entities/repository_container.go.tpl")
	exp.page.containerInterface = pager.Of("go-entities/repository_container_interface.go.tpl")

	return exp
}

func (e *RepositoryExporter) ExportPerFile( //nolint:funlen // not need
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.GoEntityRepositorySpec)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	pipeline, err := e.buildPipeline(ctx, params, spec)
	if err != nil {
		return nil, fmt.Errorf("failed to collect exporting data: %w", err)
	}

	pages := make([]*exporter.ExportedPage, 0, e.calculatePages(params, spec))
	repositories := make([]*Repository, 0, params.Schema.Tables.Len())

	repoNameMaxLength := 0
	repoInterfaceNameMaxLength := 0

	for _, table := range params.Schema.Tables.List() {
		entity := e.entityMapper.MapEntity(table, pipeline.packages.entity)
		repository := buildRepository(entity, pipeline.packages.repo, pipeline.packages.interfaces)

		if len(repository.Interface.Name) > repoInterfaceNameMaxLength {
			repoInterfaceNameMaxLength = len(repository.Interface.Name)
		}

		pkProps := e.propertyMapper.mapColumns(table.GetPKColumns(), nil)
		e.allocateRepositoryFilters(entity, repository, pipeline.packages.filters, pkProps)

		if len(repository.Name) > repoNameMaxLength {
			repoNameMaxLength = len(repository.Name)
		}

		repositories = append(repositories, repository)

		entityRepos := []*Repository{}
		if spec.Repositories.Interfaces.Place == config.GoEntityRepositorySpecRepoInterfacesPlaceWithEntity {
			entityRepos = []*Repository{repository}
		}

		page, eerr := e.entityGenerator.GenerateEntity(&GenerateEntityParams{
			Entity:       entity,
			Package:      pipeline.packages.entity,
			Repositories: entityRepos,
		})
		if eerr != nil {
			return nil, fmt.Errorf("failed to generate entity %q: %w", entity.Name, eerr)
		}

		pages = append(pages, page)

		repoFileName := fmt.Sprintf("%s.go", table.Name.Singular().Lower().Value)

		repoFile := golang.NewFile(repoFileName, pipeline.packages.repo)

		page, rerr := e.page.repo.Export(
			fmt.Sprintf("%s/%s", pipeline.packages.repo.ProjectRelativePath, repoFileName),
			map[string]stick.Value{
				"entityPackage": pipeline.packages.entity,
				"package":       pipeline.packages.repo,
				"_file":         repoFile,
				"schema": map[string]interface{}{
					"Repositories":               []*Repository{repository},
					"RepoNameMaxLength":          repoNameMaxLength,
					"RepoInterfaceNameMaxLength": repoInterfaceNameMaxLength,
					"GenInterfaces":              spec.Repositories.Interfaces.Place == config.GoEntityRepositorySpecRepoInterfacesPlaceWithRepository,                                             //nolint:lll // not need
					"GenFilters":                 spec.Repositories.Interfaces.Place == "" || spec.Repositories.Interfaces.Place == config.GoEntityRepositorySpecRepoInterfacesPlaceWithRepository, //nolint:lll // not need
				},
			},
		)
		if rerr != nil {
			return nil, rerr
		}
		pages = append(pages, page)

		if spec.Repositories.Interfaces.Place == config.GoEntityRepositorySpecRepoInterfacesPlaceEntity {
			entityRepoPageName := fmt.Sprintf(
				"%s/%s_repo.go",
				pipeline.packages.entity.ProjectRelativePath,
				entity.Table.Name.Singular().Lower(),
			)

			entityRepoP, ererr := e.page.entityRepo.Export(entityRepoPageName, map[string]stick.Value{
				"schema": map[string]interface{}{
					"Repositories": []*Repository{repository},
				},
				"_file": golang.NewFile(entityRepoPageName, pipeline.packages.entity),
			})
			if ererr != nil {
				return nil, ererr
			}

			pages = append(pages, entityRepoP)
		}
	}

	if spec.Repositories.Container.StructName != "" {
		contFileName := strings.ToLower(spec.Repositories.Container.StructName)

		containerPage := e.page.container
		if spec.Repositories.Interfaces.Place != "" {
			containerPage = e.page.containerInterface
		}

		containerGoFile := golang.NewFile(contFileName, pipeline.packages.repo)
		containerGoFile.ImportShared(golang.PackageFromFullName("github.com/jmoiron/sqlx"))
		containerGoFile.ImportLocal(pipeline.packages.entity)

		page, rerr := containerPage.Export(
			fmt.Sprintf("%s/%s.go", pipeline.packages.repo.ProjectRelativePath, contFileName),
			map[string]stick.Value{
				"_file": containerGoFile,
				"schema": map[string]interface{}{
					"Repositories":               repositories,
					"RepoNameMaxLength":          repoNameMaxLength,
					"RepoInterfaceNameMaxLength": repoInterfaceNameMaxLength,
					"Container": map[string]interface{}{
						"Name": ds.NewString(spec.Repositories.Container.StructName).Pascal().String(),
					},
					"GenInterfaces": spec.Repositories.Interfaces.Place == config.GoEntityRepositorySpecRepoInterfacesPlaceWithRepository, //nolint:lll // not need
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

func (e *RepositoryExporter) allocateRepositoryFilters(
	entity *Entity,
	repo *Repository,
	filtersPkg *golang.Package,
	props *goProperties,
) {
	repo.Filters.List = createRepositoryEntityFilter(entity, "List", filtersPkg, props)
	repo.Filters.Get = createRepositoryEntityFilter(entity, "Get", filtersPkg, props)
	repo.Filters.Delete = createRepositoryEntityFilter(entity, "Delete", filtersPkg, props)
}

func (e *RepositoryExporter) calculatePages(
	params *exporter.ExportParams,
	spec *config.GoEntityRepositorySpec,
) int {
	const defaultPageTypes = 2

	pageTypes := defaultPageTypes
	if spec.Repositories.Interfaces.Place == config.GoEntityRepositorySpecRepoInterfacesPlaceEntity {
		pageTypes++
	}

	pagesLen := params.Schema.Tables.Len() * pageTypes
	if spec.Repositories.Container.StructName != "" {
		pagesLen++
	}

	return pagesLen
}

func (e *RepositoryExporter) buildRepositoryPackage(
	spec *config.GoEntityRepositorySpec,
	goModule string,
) (*golang.Package, error) {
	pkgName := "repositories"
	if spec.Repositories.Package != "" {
		pkgName = spec.Repositories.Package
	}

	return golang.BuildPackage(pkgName, goModule)
}
