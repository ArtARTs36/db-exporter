package main

import (
	"context"
	"github.com/artarts36/db-exporter/internal/cli/cmd"
	"github.com/artarts36/db-exporter/internal/cli/config"
	"github.com/artarts36/db-exporter/internal/cli/mcp"
	"github.com/artarts36/db-exporter/internal/cli/mcp/transport"
	"github.com/artarts36/db-exporter/internal/cli/task"
	"os"
	"strings"
	"time"

	"github.com/artarts36/db-exporter/internal/exporter/factory"
	"github.com/artarts36/db-exporter/internal/shared/fs"
	"github.com/artarts36/db-exporter/internal/shared/git"
	"github.com/artarts36/db-exporter/internal/shared/migrations"
	"github.com/artarts36/db-exporter/internal/template"
	"github.com/artarts36/db-exporter/templates"
	"github.com/artarts36/singlecli"

	"github.com/tyler-sommer/stick"
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
		Args:   []*cli.ArgDefinition{},
		Opts: []*cli.OptDefinition{
			{
				Name:        "config",
				Description: "Path to config file (yaml), default: ./.db-exporter.yaml",
				WithValue:   true,
			},
			{
				Name:        "tasks",
				Description: "task names of config file",
				WithValue:   true,
			},
			{
				Name:        "mcp",
				Description: "run db-exporter in MCP mode. Enum: [console, http]",
				WithValue:   true,
			},
		},
		UsageExamples: []*cli.UsageExample{
			{
				Command:     "db-exporter --config db.yaml",
				Description: "Run db-exporter with custom config path",
			},
		},
	}

	application.RunWithGlobalArgs(context.Background())
}

func run(ctx *cli.Context) error {
	if ctx.HasOpt("mcp") {
		return runMCP(ctx)
	}
	return runCliApp(ctx)
}

func runMCP(ctx *cli.Context) error {
	fsystem := fs.NewLocal()

	cfg, err := loadConfig(ctx, fsystem)
	if err != nil {
		return err
	}

	var tp transport.Transport

	if mode, _ := ctx.GetOpt("mcp"); mode == "http" {
		tp = transport.NewHTTP()
	} else {
		tp = transport.NewConsole(time.Minute, os.Stdin, os.Stdout)
	}

	return mcp.Create(cfg, tp).Run()
}

func runCliApp(ctx *cli.Context) error {
	fsystem := fs.NewLocal()

	cfg, err := loadConfig(ctx, fsystem)
	if err != nil {
		return err
	}

	tasks := make([]string, 0)
	if taskNames, ok := ctx.GetOpt("tasks"); ok {
		tasks = strings.Split(taskNames, ",")
	}

	command := newCommand(ctx, fsystem)

	return command.Run(ctx.Context, &cmd.CommandRunParams{
		Config:    cfg,
		TaskNames: tasks,
	})
}

func newCommand(ctx *cli.Context, fs fs.Driver) *cmd.Command {
	renderer := createRenderer()

	return cmd.NewCommand(
		migrations.NewTableDetector(),
		task.NewExportActivityRunner(fs, renderer, factory.CreateExporters(renderer)),
		ctx.Output.PrintMarkdownTable,
		cmd.NewCommit(git.NewGit("git", git.GithubActionsAuthorFinder())),
	)
}

func loadConfig(ctx *cli.Context, fs fs.Driver) (*config.Config, error) {
	configPath, ok := ctx.GetOpt("config")
	if !ok {
		configPath = "./.db-exporter.yaml"
	}

	loader := config.NewLoader(fs, map[string]config.Parser{
		".yaml": config.YAMLParser(),
		".yml":  config.YAMLParser(),
		".json": config.JSONParser(),
	})

	return loader.Load(configPath)
}

func createRenderer() *template.Renderer {
	stringLoader := template.NewStringLoader()

	return template.NewRenderer(template.NewNamespaceFallbackLoader(
		template.NewNamespaceLoader(
			map[string]stick.Loader{
				"embed":  template.NewFSLoader(templates.FS),
				"local":  stick.NewFilesystemLoader("./"),
				"string": stringLoader,
			},
		),
		stringLoader,
	))
}
