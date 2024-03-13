package params

import (
	"time"

	"github.com/artarts36/db-exporter/internal/shared/fs"
)

type ActionParams struct {
	StartedAt      time.Time
	ExportParams   *ExportParams
	GeneratedFiles []fs.FileInfo
}
