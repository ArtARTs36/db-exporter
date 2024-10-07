package db

import (
	"context"
	"fmt"

	"github.com/artarts36/db-exporter/internal/schema"
)

func LoadSchemasForPool(ctx context.Context, pool *ConnectionPool) (map[string]*schema.Schema, error) {
	loader := NewPGLoader()
	schemas := map[string]*schema.Schema{}

	for db, conn := range pool.connections {
		var err error

		schemas[db], err = loader.Load(ctx, conn)
		if err != nil {
			return nil, fmt.Errorf("failed to load schema for database %q: %w", db, err)
		}

		schemas[db].SortByRelations()
	}

	return schemas, nil
}
