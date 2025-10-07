package schema

import (
	"fmt"

	"github.com/artarts36/gds"
)

type UniqueKey struct {
	Name         gds.String
	ColumnsNames *gds.Strings
}

func CreateUniqueKeyForColumn(col *Column) *UniqueKey {
	return &UniqueKey{
		Name:         *gds.NewString(fmt.Sprintf("%s_%s_uk", col.TableName.Value, col.Name.Value)),
		ColumnsNames: gds.NewStrings(col.Name.Value),
	}
}
