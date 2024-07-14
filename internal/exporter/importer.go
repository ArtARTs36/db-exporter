package exporter

import (
	"context"
	"errors"

	"github.com/artarts36/db-exporter/internal/schema"
)

type Importer interface {
	Import(ctx context.Context, schema *schema.Schema, params *ImportParams) ([]ImportedFile, error)
	ImportPerFile(ctx context.Context, sc *schema.Schema, params *ImportParams) ([]ImportedFile, error)
}

type ImportedFile struct {
	Name string
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
