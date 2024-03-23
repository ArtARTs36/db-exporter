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
