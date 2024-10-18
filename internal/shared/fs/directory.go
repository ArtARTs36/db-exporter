package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

type Directory struct {
	filesystem Driver
	dir        string
}

func NewDirectory(filesystem Driver, dir string) *Directory {
	return &Directory{filesystem: filesystem, dir: dir}
}

func (d *Directory) ReadFile(path string) ([]byte, error) {
	return d.filesystem.ReadFile(fmt.Sprintf("%s/%s", d.dir, path))
}

func (d *Directory) Up() (*Directory, error) {
	currDir := d.dir
	if currDir[len(currDir)-1] == os.PathSeparator {
		currDir = currDir[0 : len(currDir)-2]
	}

	parent := filepath.Dir(currDir)

	return &Directory{
		filesystem: d.filesystem,
		dir:        parent,
	}, nil
}
