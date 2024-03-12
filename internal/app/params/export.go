package params

type ExportParams struct {
	DriverName string
	DSN        string
	Format     string
	OutDir     string

	TablePerFile           bool
	WithDiagram            bool
	WithoutMigrationsTable bool
	Tables                 []string
	Package                string
	FilePrefix             string
	CommitMessage          string
	CommitPush             bool
}