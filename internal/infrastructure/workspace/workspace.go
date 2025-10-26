package workspace

import (
	"context"
	"github.com/artarts36/db-exporter/internal/shared/iox"
)

type Workspace interface {
	// Write file to workspace.
	Write(ctx context.Context, filename string, writer func(buffer iox.Writer) error) error
}
