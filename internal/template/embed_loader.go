package template

import (
	"embed"
	"io"

	"github.com/tyler-sommer/stick"
)

type EmbedLoader struct {
	fs embed.FS
}

type fileTemplate struct {
	name   string
	reader io.Reader
}

// NewEmbedLoader creates a new EmbedLoader with the specified root directory.
func NewEmbedLoader(fs embed.FS) *EmbedLoader {
	return &EmbedLoader{fs: fs}
}

func (t *fileTemplate) Name() string {
	return t.name
}

func (t *fileTemplate) Contents() io.Reader {
	return t.reader
}

func (l *EmbedLoader) Load(name string) (stick.Template, error) {
	f, err := l.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return &fileTemplate{name, f}, nil
}
