package proto

import "strings"

type stringsBuffer interface {
	WriteString(s string)
}

type stringsBuff struct {
	b strings.Builder
}

func (s *stringsBuff) WriteString(val string) {
	s.b.WriteString(val)
}

func (s *stringsBuff) String() string {
	return s.b.String()
}
