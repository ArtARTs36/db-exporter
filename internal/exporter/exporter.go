package exporter

import (
	"context"
	"github.com/artarts36/db-exporter/internal/schema"
)

type Exporter interface {
	Importer
	ExportPerFile(_ context.Context, sc *schema.Schema, params *ExportParams) ([]*ExportedPage, error)
	Export(ctx context.Context, schema *schema.Schema, params *ExportParams) ([]*ExportedPage, error)
}

type ExportParams struct {
	WithDiagram            bool
	WithoutMigrationsTable bool
	Package                string
}

type ExportedPage struct {
	FileName string
	Content  []byte
}

func (p *ExportedPage) Valid() bool {
	return p != nil
}
