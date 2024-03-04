package schema

import "context"

type Loader interface {
	Load(ctx context.Context, dsn string) (*Schema, error)
}
