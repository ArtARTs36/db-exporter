package main

import (
	"context"
	"github.com/artarts36/db-exporter/internal/exporter"
	"time"

	"github.com/artarts36/singlecli"

	"github.com/artarts36/db-exporter/internal/app"
)

func main() {
	application := cli.App{
		BuildInfo: &cli.BuildInfo{
			Name:      "db-exporter",
			Version:   "0.1.0",
			BuildDate: time.Now().String(),
		},
		Action: run,
		Args: []*cli.ArgDefinition{
			{
				Name:        "driver-name",
				Description: "database driver name",
				Required:    true,
				ValuesEnum: []string{
					"pg",
					"postgres",
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
				Description: "Export with diagram (only markdown)",
			},
			{
				Name:        "without-migrations-table",
				Description: "Export without migrations table",
			},
		},
		UsageExamples: []*cli.UsageExample{
			{
				Command:     "db-exporter pg \"host=postgres user=root password=root dbname=cars\" md ./docs",
				Description: "Export from postgres to markdown",
			},
		},
	}

	application.RunWithGlobalArgs(context.Background())
}

func run(ctx *cli.Context) error {
	cmd := &app.ExportCmd{}

	return cmd.Export(ctx.Context, &app.ExportParams{
		DriverName:             ctx.GetArg("driver-name"),
		DSN:                    ctx.GetArg("dsn"),
		Format:                 ctx.GetArg("format"),
		OutDir:                 ctx.GetArg("out-dir"),
		TablePerFile:           ctx.HasOpt("table-per-file"),
		WithDiagram:            ctx.HasOpt("with-diagram"),
		WithoutMigrationsTable: ctx.HasOpt("without-migrations-table"),
	})
}
