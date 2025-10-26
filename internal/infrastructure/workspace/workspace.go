package workspace

import (
	"context"

	"github.com/artarts36/db-exporter/internal/shared/indentx"
)

type Workspace interface {
	// Write file to workspace.
	Write(ctx context.Context, filename string, writer func(buffer Buffer) error) error
}

type Buffer interface {
	WriteString(s string)
	WriteIndent(ind *indentx.Indent)
	Write(p []byte) (n int, err error)
}
