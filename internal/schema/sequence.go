package schema

import "fmt"

type Sequence struct {
	Name        string `db:"name"`
	DataType    Type   `db:"-"`
	DataTypeRaw string `db:"data_type_raw"`
	Used        int    `db:"-"`
}

func (s *Sequence) Inc() {
	s.Used++
}

func CreateSequenceForColumn(col *Column) *Sequence {
	return &Sequence{
		Name:     fmt.Sprintf("%s_%s_seq", col.TableName.Value, col.Name.Value),
		DataType: col.Type,
		Used:     0,
	}
}
