package schema

import (
	"context"
	"fmt"

	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/infrastructure/conn"
	"github.com/artarts36/db-exporter/internal/schema"
)

var loaders = map[config.DatabaseDriver]Loader{
	config.DatabaseDriverPostgres: NewPGLoader(),
	config.DatabaseDriverDBML:     NewDBMLLoader(),
	config.DatabaseDriverMySQL:    NewMySQLLoader(),
}

func LoadForPool(ctx context.Context, pool *conn.Pool) (map[string]*schema.Schema, error) {
	schemas := map[string]*schema.Schema{}

	for db, con := range pool.All() {
		var err error

		loader, ok := loaders[con.Database().Driver]
		if !ok {
			return nil, fmt.Errorf("schema loader for driver %q not found", con.Database().Driver)
		}

		schemas[db], err = loader.Load(ctx, con)
		if err != nil {
			return nil, fmt.Errorf("failed to load schema for database %q: %w", db, err)
		}

		schemas[db].SortByRelations()
	}

	return schemas, nil
}
