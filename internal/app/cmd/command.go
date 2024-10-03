package cmd

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/db"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/fs"
	"github.com/artarts36/db-exporter/internal/shared/migrations"
	"log/slog"
)

type Command struct {
	migrationsTblDetector *migrations.TableDetector

	exportRunner *ExportRunner

	tablePrinter tablePrinter
}

type tablePrinter func(headers []string, rows [][]string)

func NewCommand(
	migrationsTblDetector *migrations.TableDetector,
	exportRunner *ExportRunner,
	tblPrinter tablePrinter,
) *Command {
	return &Command{
		migrationsTblDetector: migrationsTblDetector,
		exportRunner:          exportRunner,
		tablePrinter:          tblPrinter,
	}
}

func (c *Command) Run(ctx context.Context, cfg *config.Config) error {
	generatedFiles, err := c.run(ctx, cfg)
	if err != nil {
		return err
	}

	if cfg.Options.PrintStat {
		c.printStat(generatedFiles)
	}

	return nil
}

func (c *Command) run(ctx context.Context, cfg *config.Config) ([]fs.FileInfo, error) {
	connections := db.NewConnectionPool()

	err := connections.Setup(cfg.Databases)
	if err != nil {
		return nil, fmt.Errorf("failed to setup database connection pool: %w", err)
	}

	defer func() {
		closeErr := connections.Close()

		if closeErr == nil {
			slog.InfoContext(ctx, "[command] db connections closed")

			return
		}

		slog.
			With(slog.Any("err", err)).
			ErrorContext(ctx, "[command] failed to close db connections")
	}()

	schemas, err := db.LoadSchemasForPool(ctx, connections)
	if err != nil {
		return nil, err
	}

	if !cfg.Options.WithMigrationsTable {
		for _, sc := range schemas {
			sc.Tables = sc.Tables.Reject(func(table *schema.Table) bool {
				return c.migrationsTblDetector.IsMigrationsTable(table.Name.Value, table.ColumnsNames())
			})
		}
	}

	generatedFiles := make([]fs.FileInfo, 0)

	for _, task := range cfg.Tasks {
		for _, activity := range task {
			genFiles, genErr := c.exportRunner.Run(ctx, &RunExportParams{
				Activity: activity,
				Schema:   schemas[activity.Database],
			})
			if err != nil {
				return nil, genErr
			}

			generatedFiles = append(generatedFiles, genFiles...)
		}
	}

	return generatedFiles, nil
}

func (c *Command) printStat(generatedFiles []fs.FileInfo) {
	rows := make([][]string, 0, len(generatedFiles))

	for _, file := range generatedFiles {
		rows = append(rows, []string{
			file.Path,
			fmt.Sprintf("%d", file.Size),
		})
	}

	c.tablePrinter([]string{"file", "size"}, rows)
}
