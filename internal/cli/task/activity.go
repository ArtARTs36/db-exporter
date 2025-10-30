package task

import (
	"context"
	"github.com/artarts36/db-exporter/internal/cli/config"
	"github.com/artarts36/db-exporter/internal/infrastructure/conn"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/fs"
)

type ActivityRunner interface {
	Run(ctx context.Context, params *ActivityRunParams) (*ActivityResult, error)
}

type ActivityResult struct {
	Export *ExportActivityResult
}

type ExportActivityResult struct {
	files []fs.FileInfo
}

type ActivityRunParams struct {
	Activity config.Activity
	Schema   *schema.Schema
	Conn     *conn.Connection
}

func NewActivityResult() *ActivityResult {
	return &ActivityResult{
		Export: &ExportActivityResult{
			files: []fs.FileInfo{},
		},
	}
}

func (r *ActivityResult) Merge(that *ActivityResult) {
	if that.Export != nil {
		r.Export.merge(that.Export)
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
