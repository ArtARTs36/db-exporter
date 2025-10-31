package goentity

import (
	"context"
	"fmt"

	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/shared/golang"
)

type erPipeline struct {
	packages struct {
		repo       *golang.Package
		entity     *golang.Package
		filters    *golang.Package
		interfaces *golang.Package
	}

	store struct {
		repoNameMaxLength      int
		repoInterfaceMaxLength int
	}
}

func (e *RepositoryExporter) buildPipeline(
	ctx context.Context,
	params *exporter.ExportParams,
	spec *EntityRepositorySpecification,
) (*erPipeline, error) {
	pipeline := &erPipeline{}

	goModule := buildGoModule(ctx, e.goModFinder, spec.GoModule, params.Directory)

	repoPkg, err := e.buildRepositoryPackage(spec, goModule)
	if err != nil {
		return nil, fmt.Errorf("failed to build repository package: %w", err)
	}
	pipeline.packages.repo = repoPkg

	entityPkg, err := buildEntityPackage(spec.Entities.Package, goModule)
	if err != nil {
		return nil, fmt.Errorf("failed to build entity package: %w", err)
	}
	pipeline.packages.entity = entityPkg

	pipeline.packages.filters = repoPkg
	pipeline.packages.interfaces = repoPkg
	if spec.Repositories.Interfaces.Place == RepositorySpecRepoInterfacesPlaceWithEntity ||
		spec.Repositories.Interfaces.Place == RepositorySpecRepoInterfacesPlaceEntity {
		pipeline.packages.filters = entityPkg
		pipeline.packages.interfaces = entityPkg
	}

	return pipeline, nil
}
