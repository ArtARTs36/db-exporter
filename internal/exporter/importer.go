package exporter

import (
	"context"
	"errors"
	"github.com/artarts36/db-exporter/internal/db"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/fs"
)

type Importer interface {
	Import(ctx context.Context, params *ImportParams) ([]ImportedFile, error)
	ImportPerFile(ctx context.Context, params *ImportParams) ([]ImportedFile, error)
}

type ImportParams struct {
	Conn        *db.Connection
	Schema      *schema.Schema
	Directory   *fs.Directory
	TableFilter func(tableName string) bool
}

type ImportedFile struct {
	AffectedRows map[string]int64
	Name         string
}

type unimplementedImporter struct{}

func (unimplementedImporter) Import(_ context.Context, _ *ImportParams) ([]ImportedFile, error) {
	return nil, errors.New("import unimplemented")
}

func (unimplementedImporter) ImportPerFile(_ context.Context, _ *ImportParams) ([]ImportedFile, error) {
	return nil, errors.New("import per file unimplemented")
}
