package golang

import (
	"errors"
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/ds"
	"github.com/artarts36/db-exporter/internal/shared/fs"
	"os"
)

type ModFinder struct {
}

func NewModFinder() *ModFinder {
	return &ModFinder{}
}

type ModFile struct {
	Module string
}

func (f *ModFinder) FindIn(dir *fs.Directory) (*ModFile, error) {
	content, err := dir.ReadFile("go.mod")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if err == nil {
		return f.parseFile(content)
	}

	parentDir, err := dir.Up()
	if err != nil {
		return nil, fmt.Errorf("failed to get parent directory: %w", err)
	}
	content, err = parentDir.ReadFile("go.mod")
	if err != nil {
		return nil, err
	}

	return f.parseFile(content)
}

func (f *ModFinder) parseFile(content []byte) (*ModFile, error) {
	module := ds.NewString(string(content)).FirstLine().TrimPrefix("module ").TrimSpaces()
	if module.IsEmpty() {
		return nil, errors.New("failed to parse go.mod: module not found")
	}

	return &ModFile{
		Module: module.Value,
	}, nil
}
