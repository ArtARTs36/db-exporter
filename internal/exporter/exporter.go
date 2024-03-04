package exporter

import (
	"context"

	"github.com/artarts36/db-exporter/internal/schema"
)

type Exporter interface {
	Export(ctx context.Context, schema *schema.Schema, params *ExportParams) ([]*ExportedPage, error)
}

type ExportParams struct {
	TablePerFile bool
}

type ExportedPage struct {
	FileName string
	Content  []byte
}
