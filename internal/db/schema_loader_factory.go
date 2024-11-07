package db

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"

	"github.com/artarts36/db-exporter/internal/schema"
)

var schemaLoaders = map[config.DatabaseDriver]SchemaLoader{
	config.DatabaseDriverPostgres: NewPGLoader(),
	config.DatabaseDriverDBML:     NewDBMLLoader(),
}

func LoadSchemasForPool(ctx context.Context, pool *ConnectionPool) (map[string]*schema.Schema, error) {
	schemas := map[string]*schema.Schema{}

	for db, conn := range pool.connections {
		var err error

		loader, ok := schemaLoaders[conn.cfg.Driver]
		if !ok {
			return nil, fmt.Errorf("schema loader for driver %q not found", conn.cfg.Driver)
		}

		schemas[db], err = loader.Load(ctx, conn)
		if err != nil {
			return nil, fmt.Errorf("failed to load schema for database %q: %w", db, err)
		}

		schemas[db].SortByRelations()
	}

	return schemas, nil
}
