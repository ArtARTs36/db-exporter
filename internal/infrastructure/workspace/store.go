package workspace

import "github.com/artarts36/db-exporter/internal/shared/fs"

type store struct {
	files []*fs.FileInfo
}

func newStore() *store {
	return &store{
		files: make([]*fs.FileInfo, 0),
	}
}

func (s *store) Add(file *fs.FileInfo) {
	s.files = append(s.files, file)
}
