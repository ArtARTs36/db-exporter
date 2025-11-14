package sqltype

import (
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/golang"
)

func MapGoType(driver schema.DatabaseDriver, typ schema.DataType) golang.Type {
	if driver == schema.DatabaseDriverPostgres {
		return mapGoTypeFromPG(typ)
	}

	if driver == schema.DatabaseDriverDBML {
		return mapGoTypeFromDBML(typ)
	}

	return golang.TypeString
}
