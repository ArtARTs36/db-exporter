package exporter

import (
	"context"
	"errors"
	"github.com/artarts36/db-exporter/internal/schema"
)

type Importer interface {
	Import(ctx context.Context, schema *schema.Schema, params *ExportParams) error
	ImportPerFile(ctx context.Context, sc *schema.Schema, params *ExportParams) error
}

type unimplementedImporter struct{}

func (unimplementedImporter) Import(_ context.Context, _ *schema.Schema, _ *ExportParams) error {
	return errors.New("import unimplemented")
}

func (unimplementedImporter) ImportPerFile(_ context.Context, _ *schema.Schema, _ *ExportParams) error {
	return errors.New("import per file unimplemented")
}
