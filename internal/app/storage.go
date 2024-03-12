package app

import (
	"fmt"
	"log"

	"github.com/artarts36/db-exporter/internal/exporter"
	"github.com/artarts36/db-exporter/internal/shared/fs"
)

type pageStorage struct {
	fs fs.Driver
}

type savePageParams struct {
	Dir        string
	FilePrefix string
}

func (s *pageStorage) Save(pages []*exporter.ExportedPage, params *savePageParams) error {
	if !s.fs.Exists(params.Dir) {
		log.Printf("[pagestorage] creating directory %q", params.Dir)

		err := s.fs.Mkdir(params.Dir)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	for _, page := range pages {
		path := s.createPath(page, params)

		log.Printf("[pagestorage] saving %q", path)

		err := s.fs.CreateFile(path, page.Content)
		if err != nil {
			return fmt.Errorf("unable to write file %q: %w", path, err)
		}
	}

	return nil
}

func (s *pageStorage) createPath(page *exporter.ExportedPage, params *savePageParams) string {
	return fmt.Sprintf("%s/%s%s", params.Dir, params.FilePrefix, page.FileName)
}
