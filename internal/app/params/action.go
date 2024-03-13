package params

import "time"

type ActionParams struct {
	StartedAt           time.Time
	ExportParams        *ExportParams
	GeneratedFilesPaths []string
}
