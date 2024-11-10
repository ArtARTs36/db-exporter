package sqltype

import (
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/golang"
)

func MapGoType(driver config.DatabaseDriver, typ schema.Type) golang.Type {
	if driver == config.DatabaseDriverPostgres {
		return mapGoTypeFromPG(typ)
	}

	if driver == config.DatabaseDriverDBML {
		return mapGoTypeFromDBML(typ)
	}

	return golang.TypeString
}
