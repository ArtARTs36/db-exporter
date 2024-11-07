package schema

import "fmt"

type Sequence struct {
	Name             string   `db:"name"`
	DataType         string   `db:"data_type"`
	PreparedDataType DataType `db:"-"`
	Used             int      `db:"-"`
}

func (s *Sequence) Inc() {
	s.Used++
}

func CreateSequenceForColumn(col *Column) *Sequence {
	return &Sequence{
		Name:             fmt.Sprintf("%s_%s_seq", col.TableName.Value, col.Name.Value),
		DataType:         col.Type.Value,
		PreparedDataType: col.PreparedType,
		Used:             0,
	}
}
