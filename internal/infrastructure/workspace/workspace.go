package workspace

import (
	"context"
	"github.com/artarts36/db-exporter/internal/shared/iox"
)

type Workspace interface {
	// Write file to workspace.
	Write(ctx context.Context, file *WritingFile) error
}

type WritingFile struct {
	Filename string
	Writer   func(buffer iox.Writer) error
	Indent   *iox.Indent
}
