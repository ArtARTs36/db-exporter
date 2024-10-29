package schema

type Sequence struct {
	Name             string   `db:"name"`
	DataType         string   `db:"data_type"`
	PreparedDataType DataType `db:"-"`
	Used             int      `db:"-"`
}
