package main

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/app/actions"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter"
	"github.com/artarts36/db-exporter/internal/shared/migrations"
	"github.com/artarts36/db-exporter/internal/template"
	"github.com/artarts36/db-exporter/templates"
	"github.com/artarts36/singlecli"
	"github.com/tyler-sommer/stick"
	"log/slog"

	"github.com/artarts36/db-exporter/internal/app/cmd"
	"github.com/artarts36/db-exporter/internal/shared/fs"
)

var (
	Version   = "0.1.0"
	BuildDate = "2024-03-09 03:09:15"
)

func main() { //nolint:funlen // not need
	application := cli.App{
		BuildInfo: &cli.BuildInfo{
			Name:      "db-exporter",
			Version:   Version,
			BuildDate: BuildDate,
		},
		Action: run,
		Args:   []*cli.ArgDefinition{},
		Opts: []*cli.OptDefinition{
			{
				Name:        "config",
				Description: "Path to config file (yaml)",
			},
			{
				Name:        "tasks",
				Description: "task names of config file",
			},
		},
		UsageExamples: []*cli.UsageExample{
			{
				Command:     "db-exporter --config db.yaml",
				Description: "Commit db-exporter with custom config path",
			},
		},
	}

	application.RunWithGlobalArgs(context.Background())
}

func run(ctx *cli.Context) error {
	fsystem := fs.NewLocal()

	cfg, err := loadConfig(ctx)
	if err != nil {
		return err
	}

	command := newCommand(ctx, fsystem)

	return command.Run(ctx.Context, cfg)
}

func newCommand(ctx *cli.Context, fs fs.Driver) *cmd.Command {
	renderer := createRenderer(fs)

	return cmd.NewCommand(
		migrations.NewTableDetector(),
		cmd.NewExportRunner(fs, map[string]actions.Action{}, renderer, exporter.CreateExporters(renderer)),
		ctx.Output.PrintMarkdownTable,
	)
}

func loadConfig(ctx *cli.Context) (*config.Config, error) {
	configPath, ok := ctx.GetOpt("config")
	if !ok {
		configPath = "./db-exporter.yaml"
	}

	loader := &config.Loader{}

	return loader.Load(configPath)
}

func createRenderer(fs fs.Driver) *template.Renderer {
	const localTemplatesFolder = "./db-exporter-templates"

	var templateLoader stick.Loader

	if fs.Exists(localTemplatesFolder) {
		slog.Debug(fmt.Sprintf("[main] loading templates from folder %q", localTemplatesFolder))

		templateLoader = stick.NewFilesystemLoader(localTemplatesFolder)
	} else {
		slog.Debug("[main] loading templates from embedded files")

		templateLoader = template.NewEmbedLoader(templates.FS)
	}

	return template.NewRenderer(templateLoader)
}
