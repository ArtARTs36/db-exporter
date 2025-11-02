package iox

import "strings"

type Writer interface {
	WriteString(s string)
	WriteInline(s string)
	WriteNewLine()
	String() string
	Bytes() []byte

	IncIndent() Writer
	WithoutIndent() Writer
	WriteByte(b byte)
	Len() int
	Write(p []byte) (n int, err error)
}

type sbWriter struct {
	b      *strings.Builder
	indent *Indent
}

func NewWriterWithIndent(indent *Indent) Writer {
	return &sbWriter{b: &strings.Builder{}, indent: indent}
}

func NewWriter() Writer {
	return NewWriterWithIndent(zeroIndent)
}

func (s *sbWriter) WriteString(val string) {
	s.b.WriteString(s.indent.Curr())
	s.b.WriteString(val)
}

func (s *sbWriter) WriteInline(val string) {
	s.b.WriteString(val)
}

func (s *sbWriter) String() string {
	return s.b.String()
}

func (s *sbWriter) IncIndent() Writer {
	return &sbWriter{b: s.b, indent: s.indent.Next()}
}

func (s *sbWriter) WriteNewLine() {
	s.b.WriteString("\n")
}

func (s *sbWriter) WithoutIndent() Writer {
	return &sbWriter{b: s.b, indent: ZeroIndent()}
}

func (s *sbWriter) Bytes() []byte {
	return []byte(s.b.String())
}

func (s *sbWriter) WriteByte(b byte) {
	s.b.WriteByte(b)
}

func (s *sbWriter) Len() int {
	return s.b.Len()
}

func (s *sbWriter) Write(p []byte) (n int, err error) {
	return s.b.Write(p)
}
