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

func (s *Service) write(buf stringsBuffer, indent *iox.Indent) {
	if s.CommentTop != "" {
		buf.WriteString("// " + s.CommentTop + "\n")
	}

	buf.WriteString("service " + s.Name + " {")

	if len(s.Procedures) == 0 {
		buf.WriteString("}" + "\n")
		return
	}

	for _, procedure := range s.Procedures {
		buf.WriteString("\n")
		procedure.write(buf, indent.Next())
	}

	buf.WriteString("}" + "\n")
}

func (s *ServiceProcedure) write(buf stringsBuffer, indent *iox.Indent) {
	if s.CommentTop != "" {
		buf.WriteString(indent.Curr())
		buf.WriteString("// " + s.CommentTop + "\n")
	}

	buf.WriteString(indent.Curr())
	buf.WriteString("rpc " + s.Name + "(" + s.Param + ") returns (" + s.Returns + ") {")

	if len(s.Options) == 0 {
		buf.WriteString("};\n")
		return
	}

	for _, option := range s.Options {
		buf.WriteString("\n")
		option.write(buf, indent.Next())
	}

	buf.WriteString(indent.Curr())
	buf.WriteString("}\n")
}

func (opt *ServiceProcedureOption) write(buf stringsBuffer, indent *iox.Indent) {
	buf.WriteString(indent.Curr())
	buf.WriteString("option (" + opt.Name + ") = {")

	if len(opt.Params) > 0 {
		buf.WriteString("\n")

		paramsIndent := indent.Next()

		for k, v := range opt.Params {
			buf.WriteString(paramsIndent.Curr())
			buf.WriteString("" + k + ": ")
			buf.WriteString(strconv.Quote(v))
		}

		buf.WriteString("\n")
	}

	buf.WriteString(indent.Curr())
	buf.WriteString("};\n")
}
