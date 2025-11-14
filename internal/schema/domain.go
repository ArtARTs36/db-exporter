package schema

type Domain struct {
	Name     string
	DataType DataType

	ConstraintName string
	CheckClause    string

	Used int
}
