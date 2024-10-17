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
	propertyMapper  *GoPropertyMapper
}

func NewRepositoryExporter(
	pager *common.Pager,
	goModFinder *golang.ModFinder,
	entityMapper *EntityMapper,
	entityGenerator *EntityGenerator,
	propertyMapper *GoPropertyMapper,
) *RepositoryExporter {
	return &RepositoryExporter{
		pager:           pager,
		goModFinder:     goModFinder,
		entityMapper:    entityMapper,
		entityGenerator: entityGenerator,
		propertyMapper:  propertyMapper,
	}
}

func (e *RepositoryExporter) ExportPerFile( //nolint:funlen // not need
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	const pageTypes = 2

	spec, ok := params.Spec.(*config.GoEntityRepositorySpec)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	goModule := buildGoModule(ctx, e.goModFinder, spec.GoModule, params.Directory)

	repoPkg, err := e.buildRepositoryPackage(spec, goModule)
	if err != nil {
		return nil, fmt.Errorf("failed to build repository package: %w", err)
	}

	entityPkg, err := buildEntityPackage(spec.Entities.Package, goModule)
	if err != nil {
		return nil, fmt.Errorf("failed to build entity package: %w", err)
	}

	pages := make([]*exporter.ExportedPage, 0, e.calculatePages(params, spec))
	repositories := make([]*Repository, 0, params.Schema.Tables.Len())

	repoPage := e.pager.Of("go-entities/repository.go.tpl")
	repoNameMaxLength := 0
	repoInterfaceNameMaxLength := 0

	filtersPkg := repoPkg
	repoInterfacePkg := repoPkg
	if spec.Repositories.Interfaces.Place == config.GoEntityRepositorySpecRepoInterfacesPlaceWithEntity || spec.Repositories.Interfaces.Place == config.GoEntityRepositorySpecRepoInterfacesPlaceEntity {
		filtersPkg = entityPkg
		repoInterfacePkg = entityPkg
	}

	entityRepoPage := e.pager.Of("go-entities/entity_repos.go.tpl")

	for _, table := range params.Schema.Tables.List() {
		entity := e.entityMapper.MapEntity(table, entityPkg)
		repository := buildRepository(entity, repoPkg, repoInterfacePkg)

		if len(repository.Interface.Name) > repoInterfaceNameMaxLength {
			repoInterfaceNameMaxLength = len(repository.Interface.Name)
		}

		pkProps := e.propertyMapper.mapColumns(table.GetPKColumns(), nil)
		e.allocateRepositoryFilters(entity, repository, filtersPkg, pkProps)

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
			Package:      entityPkg,
			Repositories: entityRepos,
		})
		if eerr != nil {
			return nil, fmt.Errorf("failed to generate entity %q: %w", entity.Name, eerr)
		}

		pages = append(pages, page)

		repoFileName := fmt.Sprintf("%s.go", table.Name.Singular().Lower().Value)

		repoFile := golang.NewFile(repoFileName, repoPkg)

		page, rerr := repoPage.Export(
			fmt.Sprintf("%s/%s", repoPkg.ProjectRelativePath, repoFileName),
			map[string]stick.Value{
				"entityPackage": entityPkg,
				"package":       repoPkg,
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
			entityRepoPageName := fmt.Sprintf("%s/%s_repo.go", entityPkg.ProjectRelativePath, entity.Table.Name.Singular().Lower())

			entityRepoP, ererr := entityRepoPage.Export(entityRepoPageName, map[string]stick.Value{
				"schema": map[string]interface{}{
					"Repositories": []*Repository{
						repository,
					},
				},
				"_file": golang.NewFile(entityRepoPageName, entityPkg),
			})
			if ererr != nil {
				return nil, ererr
			}

			pages = append(pages, entityRepoP)
		}
	}

	if spec.Repositories.Container.StructName != "" {
		contFileName := strings.ToLower(spec.Repositories.Container.StructName)

		containerTpl := "go-entities/repository_container.go.tpl"
		if spec.Repositories.Interfaces.Place != "" {
			containerTpl = "go-entities/repository_container_interface.go.tpl"
		}

		containerGoFile := golang.NewFile(contFileName, repoPkg)
		containerGoFile.ImportShared(golang.PackageFromFullName("github.com/jmoiron/sqlx"))
		containerGoFile.ImportLocal(entityPkg)

		page, rerr := e.pager.Of(containerTpl).Export(
			fmt.Sprintf("%s/%s.go", repoPkg.ProjectRelativePath, contFileName),
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
