package cmd

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/db"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/fs"
	"github.com/artarts36/db-exporter/internal/shared/migrations"
	"github.com/artarts36/db-exporter/internal/task"
	"log/slog"
)

type Command struct {
	migrationsTblDetector *migrations.TableDetector

	activityRunner task.ActivityRunner

	tablePrinter tablePrinter
	committer    *Committer
}

type tablePrinter func(headers []string, rows [][]string)

func NewCommand(
	migrationsTblDetector *migrations.TableDetector,
	activityRunner task.ActivityRunner,
	tblPrinter tablePrinter,
	committer *Committer,
) *Command {
	return &Command{
		migrationsTblDetector: migrationsTblDetector,
		activityRunner:        activityRunner,
		tablePrinter:          tblPrinter,
		committer:             committer,
	}
}

type CommandRunParams struct {
	Config *config.Config
	Tasks  []string
}

func (c *Command) Run(ctx context.Context, params *CommandRunParams) error {
	if len(params.Tasks) > 0 {
		err := params.Config.OnlyTasks(params.Tasks)
		if err != nil {
			return err
		}
	}

	result, err := c.run(ctx, params.Config)
	if err != nil {
		return err
	}

	if params.Config.Options.PrintStat {
		c.printStat(result)
	}

	return nil
}

func (c *Command) run(ctx context.Context, cfg *config.Config) (*task.ActivityResult, error) {
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

	result := task.NewActivityResult()

	for _, ttask := range cfg.Tasks {
		exportGenFiles := make([]fs.FileInfo, 0)

		for _, activity := range ttask.Activities {
			conn, ok := connections.Get(activity.Database)
			if !ok {
				return nil, fmt.Errorf("failed to get connection for database %q", activity.Database)
			}

			activityResult, genErr := c.activityRunner.Run(ctx, &task.ActivityRunParams{
				Activity: activity,
				Schema:   schemas[activity.Database],
				Conn:     conn,
			})
			if genErr != nil {
				return nil, genErr
			}
			if activityResult == nil {
				return nil, fmt.Errorf("activity runner returns nil result")
			}

			result.Merge(activityResult)
		}

		if ttask.Commit.Valid() {
			err = c.committer.Commit(ctx, commitParams{
				Commit:         ttask.Commit,
				GeneratedFiles: exportGenFiles,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to commit: %w", err)
			}
		}
	}

	return result, nil
}

func (c *Command) printStat(result *task.ActivityResult) {
	printExport := func() {
		rows := make([][]string, 0, len(result.Export.GetFiles()))

		for _, file := range result.Export.GetFiles() {
			rows = append(rows, []string{
				file.Path,
				fmt.Sprintf("%d", file.Size),
			})
		}

		c.tablePrinter([]string{"file", "size"}, rows)
	}

	printImport := func() {
		countsList := make([][]string, 0, len(result.Import.GetTableRowCountMap()))
		for table, count := range result.Import.GetTableRowCountMap() {
			countsList = append(countsList, []string{
				table,
				fmt.Sprintf("%d", count),
			})
		}

		c.tablePrinter(
			[]string{"Table", "Affected Rows"},
			countsList,
		)
	}

	printExport()
	printImport()
}
