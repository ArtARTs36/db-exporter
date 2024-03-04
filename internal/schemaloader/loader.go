package schemaloader

import (
	"context"

	"github.com/artarts36/db-exporter/internal/schema"
)

type Loader interface {
	Load(ctx context.Context, dsn string) (*schema.Schema, error)
}
