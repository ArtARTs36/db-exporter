package schema

import (
	"context"

	"github.com/artarts36/db-exporter/internal/infrastructure/conn"
	"github.com/artarts36/db-exporter/internal/schema"
)

type Loader interface {
	// Load database schema
	Load(ctx context.Context, conn *conn.Connection) (*schema.Schema, error)
}
