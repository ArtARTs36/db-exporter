package iox

import "strings"

type Writer interface {
	WriteString(s string)
	WriteIndent(indent *Indent)
	Write(p []byte) (n int, err error)
	WriteByte(b byte)
	String() string
}

type fileBuffer struct {
	buf strings.Builder
}

func NewBuffer() Writer {
	return &fileBuffer{}
}

func (b *fileBuffer) WriteString(s string) {
	b.buf.WriteString(s)
}

func (b *fileBuffer) WriteIndent(ind *Indent) {
	b.buf.WriteString(ind.Curr())
}

func (b *fileBuffer) Write(p []byte) (n int, err error) {
	return b.buf.Write(p)
}

func (b *fileBuffer) String() string {
	return b.buf.String()
}

func (b *fileBuffer) WriteByte(c byte) {
	b.buf.WriteByte(c)
}
