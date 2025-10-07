package schema

import (
	"fmt"

	"github.com/artarts36/gds"
)

type PrimaryKey struct {
	Name         gds.String
	ColumnsNames *gds.Strings
}

func CreatePrimaryKeyForColumn(col *Column) *PrimaryKey {
	return &PrimaryKey{
		Name:         *gds.NewString(fmt.Sprintf("%s_%s_pk", col.TableName.Value, col.Name.Value)),
		ColumnsNames: gds.NewStrings(col.Name.Value),
	}
}
