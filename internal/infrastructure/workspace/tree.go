package workspace

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/fs"
	"log/slog"
)

type Tree interface {
	Files() []*fs.FileInfo
	Create(config Config) (Workspace, error)
}

type tree struct {
	fsDriver fs.Driver
	store    *store
}

func NewFSTree(fsDriver fs.Driver) Tree {
	return &tree{
		fsDriver: fsDriver,
		store:    newStore(),
	}
}

func (t *tree) Files() []*fs.FileInfo {
	return t.store.files
}

func (t *tree) Create(config Config) (Workspace, error) {
	if err := t.createDirectory(context.Background(), config); err != nil {
		return nil, err
	}

	return newFSWorkspace(config, t.fsDriver, t.store), nil
}

func (t *tree) createDirectory(ctx context.Context, config Config) error {
	if t.fsDriver.Exists(config.Directory) {
		return nil
	}

	slog.WarnContext(ctx, fmt.Sprintf("[pagestorage] creating directory %q", config.Directory))

	if err := t.fsDriver.Mkdir(config.Directory); err != nil {
		return fmt.Errorf("create directory %q: %w", config.Directory, err)
	}

	return nil
}
