package proto

import "strconv"

type Service struct {
	Name       string
	Procedures []*ServiceProcedure
}

type ServiceProcedure struct {
	Name    string
	Param   string
	Returns string
	Options []*ServiceProcedureOption
}

type ServiceProcedureOption struct {
	Name   string
	Params map[string]interface{}
}

func (s *Service) write(buf stringsBuffer, indent *Indent) {
	buf.WriteString("service " + s.Name + " {\n")

	for _, procedure := range s.Procedures {
		procedure.write(buf, indent.Next())
	}

	buf.WriteString("}" + "\n")
}

func (s *ServiceProcedure) write(buf stringsBuffer, indent *Indent) {
	buf.WriteString("  rpc " + s.Name + "(" + s.Param + ") returns (" + s.Returns + ") {")

	if len(s.Options) == 0 {
		buf.WriteString("};\n")
		return
	}

	for _, option := range s.Options {
		buf.WriteString("\n")
		option.write(buf, indent.Next())
	}

	buf.WriteString(indent.curr)
	buf.WriteString("}\n")
}

func (opt *ServiceProcedureOption) write(buf stringsBuffer, indent *Indent) {
	buf.WriteString(indent.curr)
	buf.WriteString("option (" + opt.Name + ") = {")

	if len(opt.Params) > 0 {
		buf.WriteString("\n")

		paramsIndent := indent.Next()

		for k, v := range opt.Params {
			buf.WriteString(paramsIndent.curr)
			buf.WriteString("" + k + ": ")
			buf.WriteString(strconv.Quote(v.(string)))
		}

		buf.WriteString("\n")
	}

	buf.WriteString(indent.curr)
	buf.WriteString("};\n")
}
