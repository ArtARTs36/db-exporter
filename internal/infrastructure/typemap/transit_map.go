package typemap

import (
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/shared/dbml"
	"github.com/artarts36/db-exporter/internal/shared/mysql"
	"github.com/artarts36/db-exporter/internal/shared/pg"
)

var transitSQLTypeMap = map[config.DatabaseDriver]map[config.DatabaseDriver]map[string]string{ //nolint:exhaustive,lll // not need
	config.DatabaseDriverDBML: {
		config.DatabaseDriverPostgres: map[string]string{
			mysql.TypeChar:    pg.TypeCharacter,
			mysql.TypeVarchar: pg.TypeCharacterVarying,
			mysql.TypeBinary:  pg.TypeBytea,
			mysql.TypeText:    pg.TypeText,

			mysql.TypeInt:     pg.TypeInteger,
			mysql.TypeInteger: pg.TypeInteger,

			mysql.TypeTimestamp: pg.TypeTimestampWithTZ,

			dbml.TypeUUID: pg.TypeUUID,
		},
	},
	config.DatabaseDriverPostgres: {
		config.DatabaseDriverDBML: {
			pg.TypeCharacter:        mysql.TypeChar,
			pg.TypeCharacterVarying: mysql.TypeVarchar,
			pg.TypeBytea:            mysql.TypeBinary,
			pg.TypeText:             mysql.TypeText,

			pg.TypeInteger:            mysql.TypeInteger,
			pg.TypeTimestampWithTZ:    mysql.TypeTimestamp,
			pg.TypeTimestampWithoutTZ: mysql.TypeTimestamp,

			pg.TypeUUID: dbml.TypeUUID,
		},
	},
}
