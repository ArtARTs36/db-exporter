package iox

import "strings"

type LineWriter struct {
	b *strings.Builder
}

func (w *LineWriter) WriteString(s string) {
	w.b.WriteString(s)
}
