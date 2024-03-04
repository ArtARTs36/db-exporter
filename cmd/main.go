package main

import (
	"context"
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
				Name:        "dsn",
				Description: "data source name",
				Required:    true,
			},
			{
				Name:        "format",
				Description: "exporting format",
				Required:    true,
				ValuesEnum: []string{
					"md",
				},
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
		},
		UsageExamples: []*cli.UsageExample{
			{
				Command:     "db-exporter \"host=postgres user=root password=root dbname=cars\" md ./docs",
				Description: "Export from postgres to markdown",
			},
		},
	}

	application.RunWithGlobalArgs(context.Background())
}

func run(ctx *cli.Context) error {
	cmd := &app.ExportCmd{}

	return cmd.Export(ctx.Context, &app.ExportParams{
		DSN:          ctx.GetArg("dsn"),
		Format:       ctx.GetArg("format"),
		OutDir:       ctx.GetArg("out-dir"),
		TablePerFile: ctx.HasOpt("table-per-file"),
	})
}
