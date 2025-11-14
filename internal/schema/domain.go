package schema

type Domain struct {
	Name     string
	DataType Type

	ConstraintName string
	CheckClause    string

	Used int
}
