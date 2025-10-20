package template

import (
	"io"
	"io/fs"

	"github.com/tyler-sommer/stick"
)

type FSLoader struct {
	fs fs.FS
}

type fileTemplate struct {
	name   string
	reader io.Reader
}

// NewFSLoader creates a new FSLoader.
func NewFSLoader(fs fs.FS) *FSLoader {
	return &FSLoader{fs: fs}
}

func (t *fileTemplate) Name() string {
	return t.name
}

func (t *fileTemplate) Contents() io.Reader {
	return t.reader
}

func (l *FSLoader) Load(name string) (stick.Template, error) {
	f, err := l.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return &fileTemplate{name, f}, nil
}
