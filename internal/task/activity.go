package task

import (
	"context"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/infrastructure/conn"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/fs"
)

type ActivityRunner interface {
	Run(ctx context.Context, params *ActivityRunParams) (*ActivityResult, error)
}

type ActivityResult struct {
	Import *ImportActivityResult
	Export *ExportActivityResult
}

type ImportActivityResult struct {
	files            []exporter.ImportedFile
	tableRowCountMap map[string]int64
}

type ExportActivityResult struct {
	files []fs.FileInfo
}

type ActivityRunParams struct {
	Activity config.Activity
	Schema   *schema.Schema
	Conn     *conn.Connection
}

type CompositeActivityRunner struct {
	exportRunner *ExportActivityRunner
	importRunner *ImportActivityRunner
}

func NewCompositeActivityRunner(
	exportRunner *ExportActivityRunner,
	importRunner *ImportActivityRunner,
) *CompositeActivityRunner {
	return &CompositeActivityRunner{
		exportRunner: exportRunner,
		importRunner: importRunner,
	}
}

func NewActivityResult() *ActivityResult {
	return &ActivityResult{
		Export: &ExportActivityResult{
			files: []fs.FileInfo{},
		},
		Import: &ImportActivityResult{
			files:            []exporter.ImportedFile{},
			tableRowCountMap: map[string]int64{},
		},
	}
}

func (r *CompositeActivityRunner) Run(ctx context.Context, expParams *ActivityRunParams) (*ActivityResult, error) {
	if expParams.Activity.IsExport() {
		return r.exportRunner.Run(ctx, expParams)
	}

	return r.importRunner.Run(ctx, expParams)
}

func (r *ActivityResult) Merge(that *ActivityResult) {
	if that.Export != nil {
		r.Export.merge(that.Export)
	}

	if that.Import != nil {
		r.Import.merge(that.Import)
	}
}

func (r *ExportActivityResult) GetFiles() []fs.FileInfo {
	if r == nil {
		return []fs.FileInfo{}
	}

	return r.files
}

func (r *ExportActivityResult) merge(that *ExportActivityResult) {
	r.files = append(r.files, that.GetFiles()...)
}

func (r *ImportActivityResult) merge(that *ImportActivityResult) {
	r.files = append(r.files, that.GetFiles()...)

	for table, count := range that.GetTableRowCountMap() {
		r.tableRowCountMap[table] = count
	}
}

func (r *ImportActivityResult) GetFiles() []exporter.ImportedFile {
	if r == nil {
		return []exporter.ImportedFile{}
	}

	return r.files
}

func (r *ImportActivityResult) GetTableRowCountMap() map[string]int64 {
	if r == nil {
		return map[string]int64{}
	}

	return r.tableRowCountMap
}
