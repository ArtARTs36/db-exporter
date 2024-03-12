package fs

import "time"

type Driver interface {
	Exists(path string) bool
	Mkdir(path string) error
	CreateFile(path string, content []byte) error
	Stat(path string) (*FileInfo, error)
}

type FileInfo struct {
	Path      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
