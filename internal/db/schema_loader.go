package db

import (
	"context"

	"github.com/artarts36/db-exporter/internal/schema"
)

type SchemaLoader interface {
	// Load database schema
	Load(ctx context.Context, conn *Connection) (*schema.Schema, error)
}
