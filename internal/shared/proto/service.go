package proto

import (
	"github.com/artarts36/db-exporter/internal/shared/iox"
	"strconv"
)

type Service struct {
	Name       string
	Procedures []*ServiceProcedure

	CommentTop string
}

type ServiceProcedure struct {
	Name    string
	Param   string
	Returns string
	Options []*ServiceProcedureOption

	CommentTop string
}

type ServiceProcedureOption struct {
	Name   string
	Params map[string]string
}

func (s *Service) write(buf iox.Writer) {
	if s.CommentTop != "" {
		buf.WriteString("// " + s.CommentTop + "\n")
	}

	buf.WriteString("service " + s.Name + " {")

	if len(s.Procedures) == 0 {
		buf.WriteString("}" + "\n")
		return
	}

	for _, procedure := range s.Procedures {
		buf.WriteNewLine()
		procedure.write(buf.IncIndent())
	}

	buf.WriteString("}" + "\n")
}

func (s *ServiceProcedure) write(buf iox.Writer) {
	if s.CommentTop != "" {
		buf.WriteString("// " + s.CommentTop + "\n")
	}

	buf.WriteString("rpc " + s.Name + "(" + s.Param + ") returns (" + s.Returns + ") {")

	if len(s.Options) == 0 {
		buf.WriteInline("};\n")
		return
	}

	for _, option := range s.Options {
		buf.WriteNewLine()
		option.write(buf.IncIndent())
	}

	buf.WriteString("}\n")
}

func (opt *ServiceProcedureOption) write(buf iox.Writer) {
	buf.WriteString("option (" + opt.Name + ") = {")

	if len(opt.Params) > 0 {
		buf.WriteNewLine()

		paramsBuf := buf.IncIndent()

		for k, v := range opt.Params {
			paramsBuf.WriteString("" + k + ": " + strconv.Quote(v))
		}

		buf.WriteNewLine()
	}

	buf.WriteString("};\n")
}
