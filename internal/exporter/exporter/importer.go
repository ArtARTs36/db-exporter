package exporter

import (
	"context"
	"github.com/artarts36/db-exporter/internal/infrastructure/conn"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/fs"
)

type Importer interface {
	Import(ctx context.Context, params *ImportParams) ([]ImportedFile, error)
}

type ImportParams struct {
	Conn        *conn.Connection
	Schema      *schema.Schema
	Directory   *fs.Directory
	TableFilter func(tableName string) bool
}

type ImportedFile struct {
	AffectedRows map[string]int64
	Name         string
}
