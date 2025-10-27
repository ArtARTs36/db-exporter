package iox

import "strings"

type Writer interface {
	WriteString(s string)
	WriteInline(s string)
	WriteNewLine()
	String() string

	IncIndent() Writer
	WithoutIndent() Writer
}

type sbWriter struct {
	b      *strings.Builder
	indent *Indent
}

func NewWriter(indent *Indent) Writer {
	return &sbWriter{b: &strings.Builder{}, indent: indent}
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
