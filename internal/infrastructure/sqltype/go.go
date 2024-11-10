package sqltype

import (
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/golang"
)

func MapGoType(driver config.DatabaseDriver, typ schema.Type) golang.Type {
	if driver == config.DatabaseDriverPostgres {
		return MapGoTypeFromPG(typ)
	}

	if driver == config.DatabaseDriverDBML {
		return MapGoTypeFromDBML(typ)
	}

	return golang.TypeString
}
