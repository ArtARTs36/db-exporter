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

func (s *pageStorage) Save(dir string, pages []*exporter.ExportedPage) error {
	if !s.fs.Exists(dir) {
		log.Printf("[pagestorage] creating directory %q", dir)

		err := s.fs.Mkdir(dir)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	for _, page := range pages {
		path := fmt.Sprintf("%s/%s", dir, page.FileName)

		log.Printf("[pagestorage] saving %q", path)

		err := s.fs.CreateFile(path, page.Content)
		if err != nil {
			return fmt.Errorf("unable to write file: %w", err)
		}
	}

	return nil
}
