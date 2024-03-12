package main

import (
	"context"
	"strings"

	"github.com/artarts36/singlecli"

	"github.com/artarts36/db-exporter/internal/app/actions"
	"github.com/artarts36/db-exporter/internal/app/cmd"
	"github.com/artarts36/db-exporter/internal/app/params"
	"github.com/artarts36/db-exporter/internal/exporter"
	"github.com/artarts36/db-exporter/internal/shared/fs"
	"github.com/artarts36/db-exporter/internal/shared/git"
)

var (
	Version   = "0.1.0"
	BuildDate = "2024-03-09 03:09:15"
)

func main() {
	application := cli.App{
		BuildInfo: &cli.BuildInfo{
			Name:      "db-exporter",
			Version:   Version,
			BuildDate: BuildDate,
		},
		Action: run,
		Args: []*cli.ArgDefinition{
			{
				Name:        "driver-name",
				Description: "database driver name",
				Required:    true,
				ValuesEnum: []string{
					"pg",
				},
			},
			{
				Name:        "dsn",
				Description: "data source name",
				Required:    true,
			},
			{
				Name:        "format",
				Description: "exporting format",
				Required:    true,
				ValuesEnum:  exporter.Names,
			},
			{
				Name:        "out-dir",
				Description: "Output directory",
				Required:    true,
			},
		},
		Opts: []*cli.OptDefinition{
			{
				Name:        "table-per-file",
				Description: "Export one table to one file",
			},
			{
				Name:        "with-diagram",
				Description: "Export with diagram (only md)",
			},
			{
				Name:        "without-migrations-table",
				Description: "Export without migrations table",
			},
			{
				Name:        "tables",
				Description: "Table list for export, separator: \",\"",
			},
			{
				Name:        "package",
				Description: "Package name for code gen, e.g: models",
			},
			{
				Name:        "file-prefix",
				Description: "Prefix for generated files",
			},
			{
				Name:        "commit-message",
				Description: "Add commit with generated files and your message",
			},
			{
				Name:        "commit-push",
				Description: "Push commit with generated files",
			},
		},
		UsageExamples: []*cli.UsageExample{
			{
				Command:     "db-exporter pg \"host=postgres user=root password=root dbname=cars\" md ./docs",
				Description: "Export from postgres to md",
			},
		},
	}

	application.RunWithGlobalArgs(context.Background())
}

func run(ctx *cli.Context) error {
	command := cmd.NewExportCmd(fs.NewLocal(), map[string]actions.Action{
		"commit generated files": actions.NewCommit(git.NewGit("git")),
	})

	var tables []string

	tablesOpt, hasTablesOpt := ctx.GetOpt("tables")
	if hasTablesOpt {
		tables = strings.Split(tablesOpt, ",")
	}

	pkg, _ := ctx.GetOpt("package")
	filePrefix, _ := ctx.GetOpt("file-prefix")
	commitMessage, _ := ctx.GetOpt("commit-message")

	return command.Export(ctx.Context, &params.ExportParams{
		DriverName:             ctx.GetArg("driver-name"),
		DSN:                    ctx.GetArg("dsn"),
		Format:                 ctx.GetArg("format"),
		OutDir:                 ctx.GetArg("out-dir"),
		TablePerFile:           ctx.HasOpt("table-per-file"),
		WithDiagram:            ctx.HasOpt("with-diagram"),
		WithoutMigrationsTable: ctx.HasOpt("without-migrations-table"),
		Tables:                 tables,
		Package:                pkg,
		FilePrefix:             filePrefix,
		CommitMessage:          commitMessage,
		CommitPush:             ctx.HasOpt("commit-push"),
	})
}
