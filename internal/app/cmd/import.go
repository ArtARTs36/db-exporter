package cmd

//
// import (
// 	"context"
// 	"fmt"
// 	"github.com/artarts36/db-exporter/internal/shared/ds"
// 	"log/slog"
// 	"slices"
// 	"strings"
//
// 	"github.com/tyler-sommer/stick"
//
// 	"github.com/artarts36/db-exporter/internal/app/params"
// 	"github.com/artarts36/db-exporter/internal/db"
// 	"github.com/artarts36/db-exporter/internal/exporter"
// 	"github.com/artarts36/db-exporter/internal/schema"
// 	"github.com/artarts36/db-exporter/internal/shared/fs"
// 	"github.com/artarts36/db-exporter/internal/shared/migrations"
// 	"github.com/artarts36/db-exporter/internal/template"
// 	"github.com/artarts36/db-exporter/templates"
// )
//
// type ImportCmd struct {
// 	migrationsTblDetector *migrations.TableDetector
// 	fs                    fs.Driver
// 	tablePrinter          tablePrinter
// }
//
// type tablePrinter func(headers []string, rows [][]string)
//
// func NewImportCmd(fs fs.Driver, tablePrinter tablePrinter) *ImportCmd {
// 	return &ImportCmd{
// 		migrationsTblDetector: migrations.NewTableDetector(),
// 		fs:                    fs,
// 		tablePrinter:          tablePrinter,
// 	}
// }
//
// func (a *ImportCmd) Commit(ctx context.Context, expParams *params.Config) error {
// 	driverName, err := db.CreateDriverName(expParams.DriverName)
// 	if err != nil {
// 		return err
// 	}
//
// 	connection := db.NewConnection(driverName, expParams.DSN)
// 	defer func() {
// 		closeErr := connection.Close()
//
// 		if closeErr == nil {
// 			slog.InfoContext(ctx, "[ImportCmd] db connection closed")
//
// 			return
// 		}
//
// 		slog.ErrorContext(ctx, fmt.Sprintf("failed to close db connection: %s", closeErr))
// 	}()
//
// 	loader, err := db.CreateSchemaLoader(connection)
// 	if err != nil {
// 		return fmt.Errorf("unable to create schema loader: %w", err)
// 	}
//
// 	renderer := a.createRenderer()
//
// 	exp, err := exporter.CreateExporter(expParams.Format, renderer, connection)
// 	if err != nil {
// 		return fmt.Errorf("failed to create exporter: %w", err)
// 	}
//
// 	// processing
//
// 	slog.DebugContext(ctx, fmt.Sprintf("[importcmd] loading db schema from %s", expParams.DriverName))
//
// 	sc, err := a.loadSchema(ctx, loader, expParams)
// 	if err != nil {
// 		return fmt.Errorf("unable to load schema: %w", err)
// 	}
//
// 	files, err := a.doImport(ctx, exp, sc, expParams)
// 	if err != nil {
// 		return err
// 	}
//
// 	if len(files) == 0 {
// 		slog.InfoContext(ctx, "[importcmd] no files to import")
// 	} else {
// 		filesPaths := strings.Builder{}
// 		countsMap := map[string]int64{}
//
// 		for _, file := range files {
// 			if filesPaths.Len() > 0 {
// 				filesPaths.WriteRune(',')
// 			}
//
// 			filesPaths.WriteString(file.Name)
//
// 			for table, ar := range file.AffectedRows {
// 				countsMap[table] += ar
// 			}
// 		}
//
// 		slog.InfoContext(
// 			ctx,
// 			fmt.Sprintf("[importcmd] successfully imported from %d files: %s", len(files), filesPaths.String()),
// 		)
//
// 		countsList := make([][]string, 0, len(countsMap))
// 		for table, count := range countsMap {
// 			countsList = append(countsList, []string{
// 				table,
// 				fmt.Sprintf("%d", count),
// 			})
// 		}
//
// 		a.tablePrinter(
// 			[]string{"Table", "Affected Rows"},
// 			countsList,
// 		)
// 	}
//
// 	return nil
// }
//
// func (a *ImportCmd) doImport(
// 	ctx context.Context,
// 	exp exporter.Exporter,
// 	sc *schema.Schema,
// 	params *params.Config,
// ) ([]exporter.ImportedFile, error) {
// 	var pages []exporter.ImportedFile
// 	var err error
// 	importerParams := &exporter.ImportParams{
// 		Directory: fs.NewDirectory(a.fs, params.OutDir),
// 		TableFilter: func(tableName string) bool {
// 			if len(params.Tables) > 0 && !slices.Contains(params.Tables, tableName) {
// 				return false
// 			}
//
// 			if params.WithoutMigrationsTable {
// 				table, exists := sc.Tables.Get(*ds.NewString(tableName))
// 				if exists && a.migrationsTblDetector.IsMigrationsTable(tableName, table.ColumnsNames()) {
// 					return false
// 				}
// 			}
//
// 			return true
// 		},
// 	}
//
// 	if params.TablePerFile {
// 		pages, err = exp.ImportPerFile(ctx, sc, importerParams)
// 	} else {
// 		pages, err = exp.Import(ctx, sc, importerParams)
// 	}
//
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to import: %w", err)
// 	}
//
// 	return pages, nil
// }
//
// func (a *ImportCmd) loadSchema(
// 	ctx context.Context,
// 	loader db.SchemaLoader,
// 	params *params.Config,
// ) (*schema.Schema, error) {
// 	sc, err := loader.Load(ctx, nil)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if len(params.Tables) > 0 {
// 		slog.DebugContext(ctx, "[ImportCmd] filtering tables")
//
// 		sc.Tables = sc.Tables.Only(params.Tables)
// 	}
//
// 	slog.DebugContext(ctx, "[ImportCmd] sorting tables by relations")
//
// 	sc.SortByRelations()
//
// 	if !params.WithoutMigrationsTable {
// 		return sc, nil
// 	}
//
// 	sc.Tables = sc.Tables.Reject(func(table *schema.Table) bool {
// 		return a.migrationsTblDetector.IsMigrationsTable(table.Name.Value, table.ColumnsNames())
// 	})
//
// 	return sc, nil
// }
//
// func (a *ImportCmd) createRenderer() *template.Renderer {
// 	var templateLoader stick.Loader
//
// 	if a.fs.Exists(localTemplatesFolder) {
// 		slog.Debug(fmt.Sprintf("[ImportCmd] loading templates from folder %q", localTemplatesFolder))
//
// 		templateLoader = stick.NewFilesystemLoader(localTemplatesFolder)
// 	} else {
// 		slog.Debug("[ImportCmd] loading templates from embedded files")
//
// 		templateLoader = template.NewEmbedLoader(templates.FS)
// 	}
//
// 	return template.NewRenderer(templateLoader)
// }
