package schema

type Schema struct {
	Tables map[String]*Table
}

type Table struct {
	Name    String    `db:"name"`
	Columns []*Column `db:"-"`
}

type ForeignKey struct {
	Name   String
	Table  String
	Column String
}
