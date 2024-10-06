package cmd

import (
	"context"
	"errors"
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
	Config    *config.Config
	TaskNames []string

	dbs   map[string]config.Database
	tasks map[string]config.Task
}

func (c *Command) Run(ctx context.Context, params *CommandRunParams) error {
	tasks := make(map[string]config.Task, 0)
	dbKeysSet := make(map[string]bool, 0)
	dbs := make(map[string]config.Database, 0)
	for _, name := range params.TaskNames {
		t, ok := params.Config.Tasks[name]
		if !ok {
			return fmt.Errorf("task %q not found", name)
		}
		tasks[name] = t

		for _, act := range t.Activities {
			exists := dbKeysSet[act.Database]
			if !exists {
				dbKeysSet[act.Database] = true
				dbs[act.Database] = params.Config.Databases[act.Database]
			}
		}
	}
	if len(tasks) == 0 {
		return errors.New("tasks not found")
	}

	params.tasks = tasks
	params.dbs = dbs

	result, err := c.run(ctx, params)
	if err != nil {
		return err
	}

	if params.Config.Options.PrintStat {
		c.printStat(result)
	}

	return nil
}

func (c *Command) run(ctx context.Context, params *CommandRunParams) (*task.ActivityResult, error) {
	connections := db.NewConnectionPool()

	err := connections.Setup(params.dbs)
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

	if !params.Config.Options.WithMigrationsTable {
		for _, sc := range schemas {
			sc.Tables = sc.Tables.Reject(func(table *schema.Table) bool {
				return c.migrationsTblDetector.IsMigrationsTable(table.Name.Value, table.ColumnsNames())
			})
		}
	}

	result := task.NewActivityResult()

	for taskName, ttask := range params.tasks {
		slog.InfoContext(ctx, "[command] running task", slog.String("task", taskName))

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
