package fs

import "fmt"

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
