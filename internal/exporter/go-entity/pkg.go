package goentity

import (
	"context"
	"github.com/artarts36/db-exporter/internal/shared/fs"
	"github.com/artarts36/db-exporter/internal/shared/golang"
	"log/slog"
)

func buildGoModule(ctx context.Context, finder *golang.ModFinder, mod string, dir *fs.Directory) string {
	if mod == "" {
		goMod, err := finder.FindIn(dir)
		if err != nil {
			slog.
				With(slog.Any("err", err)).
				WarnContext(ctx, "[go-entity-repository-exporter] failed to get go module")
		} else {
			return goMod.Module
		}
	}

	return mod
}

func buildEntityPackage(pkg, goModule string) (golang.Package, error) {
	if pkg == "" {
		pkg = "entities"
	}

	return golang.BuildPackage(pkg, goModule)
}
