package schema

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/infrastructure/schema/dbml"
	"github.com/artarts36/db-exporter/internal/infrastructure/schema/mysql"
	"github.com/artarts36/db-exporter/internal/infrastructure/schema/pg"

	"github.com/artarts36/db-exporter/internal/infrastructure/conn"
	"github.com/artarts36/db-exporter/internal/schema"
)

var loaders = map[schema.DatabaseDriver]Loader{
	schema.DatabaseDriverPostgres: pg.NewLoader(),
	schema.DatabaseDriverDBML:     dbml.NewLoader(),
	schema.DatabaseDriverMySQL:    mysql.NewLoader(),
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
