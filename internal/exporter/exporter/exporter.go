package exporter

import (
	"context"
	"github.com/artarts36/db-exporter/internal/infrastructure/conn"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/fs"
)

type Exporter interface {
	ExportPerFile(ctx context.Context, params *ExportParams) ([]*ExportedPage, error)
	Export(ctx context.Context, params *ExportParams) ([]*ExportedPage, error)
}

type ExportParams struct {
	Schema    *schema.Schema
	Conn      *conn.Connection
	Spec      interface{}
	Directory *fs.Directory
}

type ExportedPage struct {
	FileName string
	Content  []byte
}

func (p *ExportedPage) Valid() bool {
	return p != nil
}
