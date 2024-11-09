package schema

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/ds"
)

type PrimaryKey struct {
	Name         ds.String
	ColumnsNames *ds.Strings
}

func CreatePrimaryKeyForColumn(col *Column) *PrimaryKey {
	return &PrimaryKey{
		Name:         *ds.NewString(fmt.Sprintf("%s_%s_pk", col.TableName.Value, col.Name)),
		ColumnsNames: ds.NewStrings(col.Name.Value),
	}
}
