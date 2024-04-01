package schemaloader

import (
	"context"

	"github.com/artarts36/db-exporter/internal/schema"
)

type Loader interface {
	// Load database schema
	Load(ctx context.Context) (*schema.Schema, error)
}
