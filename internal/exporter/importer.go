package exporter

import (
	"context"
	"errors"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/fs"
)

type Importer interface {
	Import(ctx context.Context, schema *schema.Schema, params *ImportParams) ([]ImportedFile, error)
	ImportPerFile(ctx context.Context, sc *schema.Schema, params *ImportParams) ([]ImportedFile, error)
}

type ImportParams struct {
	Directory   *fs.Directory
	TableFilter func(tableName string) bool
}

type ImportedFile struct {
	AffectedRows map[string]int64
	Name         string
}

type unimplementedImporter struct{}

func (unimplementedImporter) Import(_ context.Context, _ *schema.Schema, _ *ImportParams) ([]ImportedFile, error) {
	return nil, errors.New("import unimplemented")
}

func (unimplementedImporter) ImportPerFile(
	_ context.Context,
	_ *schema.Schema,
	_ *ImportParams,
) ([]ImportedFile, error) {
	return nil, errors.New("import per file unimplemented")
}
