package workspace

import (
	"strings"

	"github.com/artarts36/db-exporter/internal/shared/indentx"
)

type fileBuffer struct {
	buf strings.Builder
}

func (b *fileBuffer) WriteString(s string) {
	b.buf.WriteString(s)
}

func (b *fileBuffer) WriteIndent(ind *indentx.Indent) {
	b.buf.WriteString(ind.Curr())
}

func (b *fileBuffer) Write(p []byte) (n int, err error) {
	return b.buf.Write(p)
}
