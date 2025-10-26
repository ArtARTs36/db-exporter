package task

import (
	"context"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/infrastructure/conn"
	"github.com/artarts36/db-exporter/internal/infrastructure/workspace"
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
	Files []fs.FileInfo
}

type ActivityRunParams struct {
	Activity      config.Activity
	Schema        *schema.Schema
	Conn          *conn.Connection
	WorkspaceTree workspace.Tree
}

func NewActivityResult() *ActivityResult {
	return &ActivityResult{
		Export: &ExportActivityResult{
			Files: []fs.FileInfo{},
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

	return r.Files
}

func (r *ExportActivityResult) merge(that *ExportActivityResult) {
	r.Files = append(r.Files, that.GetFiles()...)
}
