package cmd

import (
	"context"

	"github.com/artarts36/db-exporter/internal/app/params"
)

type Command interface {
	Run(ctx context.Context, params *params.ExportParams) error
}
