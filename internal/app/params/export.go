package params

type Config struct {
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
	CommitAuthor           string
	CommitPush             bool
	Stat                   bool
	Import                 bool
	ProtoGoPackage         string
}
