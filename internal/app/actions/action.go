package actions

import (
	"context"

	"github.com/artarts36/db-exporter/internal/app/params"
)

type Action interface {
	Supports(params *params.ActionParams) bool
	Run(ctx context.Context, params *params.ActionParams) error
}
