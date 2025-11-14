package schema

import "fmt"

type Sequence struct {
	Name        string   `db:"name"`
	DataType    DataType `db:"-"`
	DataTypeRaw string   `db:"data_type_raw"`
	Used        int      `db:"-"`
}

func (s *Sequence) Inc() {
	s.Used++
}

func (s *Sequence) UsedOnce() bool {
	return s.Used == 1
}

func CreateSequenceForColumn(col *Column) *Sequence {
	return &Sequence{
		Name:     fmt.Sprintf("%s_%s_seq", col.TableName.Value, col.Name.Value),
		DataType: col.DataType,
		Used:     0,
	}
}
