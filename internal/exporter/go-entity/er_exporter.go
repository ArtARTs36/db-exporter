package goentity

import (
	"context"
	"errors"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/golang"
	"github.com/artarts36/gds"
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
		enumString         *common.Page
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

	exp.page.repo = pager.Of("@embed/go-entities/repository.go.tpl")
	exp.page.entityRepo = pager.Of("@embed/go-entities/entity_repos.go.tpl")
	exp.page.container = pager.Of("@embed/go-entities/repository_container.go.tpl")
	exp.page.containerInterface = pager.Of("@embed/go-entities/repository_container_interface.go.tpl")
	exp.page.enumString = pager.Of("@embed/go-entities/enum_string.go.tpl")

	return exp
}

func (e *RepositoryExporter) renderEnums(pipeline *erPipeline, sch *schema.Schema) (
	[]*exporter.ExportedPage,
	map[string]*golang.StringEnum,
	error,
) {
	enums := map[string]*golang.StringEnum{}
	pages := make([]*exporter.ExportedPage, 0, len(sch.Enums))

	for _, enum := range sch.Enums {
		enumFile := golang.NewFile(fmt.Sprintf("%s.go", enum.Name.Value), pipeline.packages.entity)

		goEnum := golang.NewStringEnumOfValues(enum.Name, enum.Values)
		enums[enum.Name.Value] = goEnum

		page, enumErr := e.page.enumString.Export(
			fmt.Sprintf("%s/%s", pipeline.packages.entity.ProjectRelativePath, enumFile.Name),
			map[string]stick.Value{
				"enum":  goEnum,
				"_file": enumFile,
			},
		)
		if enumErr != nil {
			return nil, nil, fmt.Errorf("failed to generate enum %q: %w", enum.Name.Value, enumErr)
		}

		pages = append(pages, page)
	}

	return pages, enums, nil
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

	enumPages, enums, err := e.renderEnums(pipeline, params.Schema)
	if err != nil {
		return nil, fmt.Errorf("failed to render enums: %w", err)
	}
	pages = append(pages, enumPages...)

	for _, table := range params.Schema.Tables.List() {
		entity := e.entityMapper.MapEntity(&MapEntityParams{
			SourceDriver: params.Schema.Driver,
			Table:        table,
			Package:      pipeline.packages.entity,
			Enums:        enums,
		})
		repository := buildRepository(entity, pipeline.packages.repo, pipeline.packages.interfaces)

		pkProps := e.propertyMapper.mapColumns(params.Schema.Driver, table.GetPKColumns(), enums, nil)
		e.allocateRepositoryFilters(entity, repository, pipeline.packages.filters, pkProps)

		if len(repository.Name) > pipeline.store.repoNameMaxLength {
			pipeline.store.repoNameMaxLength = len(repository.Name)
		}
		if len(repository.Interface.Name) > pipeline.store.repoInterfaceMaxLength {
			pipeline.store.repoInterfaceMaxLength = len(repository.Interface.Name)
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
			WithMocks:    spec.Repositories.Interfaces.WithMocks,
		})
		if eerr != nil {
			return nil, fmt.Errorf("failed to generate entity %q: %w", entity.Name, eerr)
		}

		pages = append(pages, page)

		page, rerr := e.page.repo.Export(
			fmt.Sprintf("%s/%s", pipeline.packages.repo.ProjectRelativePath, repository.File.Name),
			map[string]stick.Value{
				"entityPackage": pipeline.packages.entity,
				"package":       pipeline.packages.repo,
				"_file":         repository.File,
				"schema": map[string]interface{}{
					"Repositories":               []*Repository{repository},
					"RepoNameMaxLength":          pipeline.store.repoNameMaxLength,
					"RepoInterfaceNameMaxLength": pipeline.store.repoInterfaceMaxLength,
					"GenInterfaces":              spec.Repositories.Interfaces.Place == config.GoEntityRepositorySpecRepoInterfacesPlaceWithRepository,                                             //nolint:lll // not need
					"GenFilters":                 spec.Repositories.Interfaces.Place == "" || spec.Repositories.Interfaces.Place == config.GoEntityRepositorySpecRepoInterfacesPlaceWithRepository, //nolint:lll // not need
					"WithMocks":                  spec.Repositories.Interfaces.WithMocks,
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
		containerGoFile.ImportLocal(pipeline.packages.entity)

		page, rerr := containerPage.Export(
			fmt.Sprintf("%s/%s.go", pipeline.packages.repo.ProjectRelativePath, contFileName),
			map[string]stick.Value{
				"_file": containerGoFile,
				"schema": map[string]interface{}{
					"Repositories":               repositories,
					"RepoNameMaxLength":          pipeline.store.repoNameMaxLength,
					"RepoInterfaceNameMaxLength": pipeline.store.repoInterfaceMaxLength,
					"Container": map[string]interface{}{
						"Name": gds.NewString(spec.Repositories.Container.StructName).Pascal().String(),
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

	pagesLen += len(params.Schema.Enums)

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
