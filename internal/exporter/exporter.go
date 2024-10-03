package exporter

import (
	"context"
	"github.com/artarts36/db-exporter/internal/db"
	"github.com/artarts36/db-exporter/internal/schema"
)

type Exporter interface {
	ExportPerFile(_ context.Context, params *ExportParams) ([]*ExportedPage, error)
	Export(ctx context.Context, params *ExportParams) ([]*ExportedPage, error)
}

type ExportParams struct {
	Schema *schema.Schema
	Conn   *db.Connection
	Spec   interface{}
}

type ExportedPage struct {
	FileName string
	Content  []byte
}

func (p *ExportedPage) Valid() bool {
	return p != nil
}
