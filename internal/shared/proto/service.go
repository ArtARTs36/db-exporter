package proto

type Service struct {
	Name       string
	Procedures []*ServiceProcedure
}

type ServiceProcedure struct {
	Name    string
	Param   string
	Returns string
}

func (s *Service) write(buf stringsBuffer) {
	buf.WriteString("service " + s.Name + " {\n")

	for _, procedure := range s.Procedures {
		procedure.write(buf)
	}

	buf.WriteString("}" + "\n")
}

func (s *ServiceProcedure) write(buf stringsBuffer) {
	buf.WriteString("    rpc " + s.Name + "(" + s.Param + ") returns (" + s.Returns + ") {}" + "\n")
}
