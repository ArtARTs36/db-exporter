package params

import (
	"github.com/artarts36/db-exporter/internal/shared/fs"
)

type ActionParams struct {
	ExportParams   *Config
	GeneratedFiles []fs.FileInfo
}
