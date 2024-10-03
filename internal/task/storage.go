package task

import (
	"fmt"
	"log/slog"

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

func (s *pageStorage) Save(pages []*exporter.ExportedPage, params *savePageParams) ([]fs.FileInfo, error) {
	writtenFiles := make([]fs.FileInfo, 0, len(pages))

	if !s.fs.Exists(params.Dir) {
		slog.Info(fmt.Sprintf("[pagestorage] creating directory %q", params.Dir))

		err := s.fs.Mkdir(params.Dir)
		if err != nil {
			return writtenFiles, fmt.Errorf("failed to create directory: %w", err)
		}
	}

	slog.Debug(fmt.Sprintf("[pagestorage] saving %d files", len(pages)))

	for _, page := range pages {
		path := s.createPath(page, params)

		slog.Debug(fmt.Sprintf("[pagestorage] saving %q", path))

		wrFile, err := s.fs.Write(path, page.Content)
		if err != nil {
			return writtenFiles, fmt.Errorf("unable to write file %q: %w", path, err)
		}

		writtenFiles = append(writtenFiles, wrFile)
	}

	slog.Info(fmt.Sprintf("[pagestorage] saved %d files", len(pages)))

	return writtenFiles, nil
}

func (s *pageStorage) createPath(page *exporter.ExportedPage, params *savePageParams) string {
	return fmt.Sprintf("%s/%s%s", params.Dir, params.FilePrefix, page.FileName)
}
