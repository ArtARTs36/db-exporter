package sqltype

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

			dbml.TypeFloat: pg.TypeDoublePrecision,
			pg.TypeFloat8:  pg.TypeFloat8,
		},
	},
	config.DatabaseDriverPostgres: {
		config.DatabaseDriverDBML: {
			pg.TypeCharacter:        mysql.TypeChar,
			pg.TypeCharacterVarying: mysql.TypeVarchar,
			pg.TypeBytea:            mysql.TypeBinary,
			pg.TypeText:             mysql.TypeText,

			pg.TypeInteger: mysql.TypeInteger,

			pg.TypeTimestampWithTZ:    mysql.TypeTimestamp,
			pg.TypeTimestampWithoutTZ: mysql.TypeTimestamp,

			pg.TypeUUID: dbml.TypeUUID,

			pg.TypeDoublePrecision: dbml.TypeFloat,
			pg.TypeFloat8:          pg.TypeFloat8,
		},

		// https://dev.mysql.com/doc/workbench/en/wb-migration-database-postgresql-typemapping.html
		config.DatabaseDriverMySQL: {
			pg.TypeInt:                mysql.TypeInt,
			pg.TypeSmallInt:           mysql.TypeSmallInt,
			pg.TypeBigint:             mysql.TypeInt,
			pg.TypeSerial:             mysql.TypeInt,
			pg.TypeSmallSerial:        mysql.TypeSmallInt,
			pg.TypeBigSerial:          mysql.TypeBigInt,
			pg.TypeBit:                mysql.TypeBit,
			pg.TypeBoolean:            mysql.TypeTinyint, // 1
			pg.TypeReal:               mysql.TypeFloat,
			pg.TypeDoublePrecision:    mysql.TypeDouble,
			pg.TypeNumeric:            mysql.TypeDecimal,
			pg.TypeDecimal:            mysql.TypeDecimal,
			pg.TypeMoney:              mysql.TypeDecimal, // 19,2
			pg.TypeCharacter:          mysql.TypeChar,
			pg.TypeChar:               mysql.TypeChar,
			pg.TypeCharacterVarying:   mysql.TypeLongText,
			pg.TypeDate:               mysql.TypeDate,
			pg.TypeTimeWithTZ:         mysql.TypeTime,
			pg.TypeTimeWithoutTZ:      mysql.TypeTime,
			pg.TypeTimestampWithTZ:    mysql.TypeDateTime,
			pg.TypeTimestampWithoutTZ: mysql.TypeDateTime,
			pg.TypeInterval:           mysql.TypeTime,
			pg.TypeBytea:              mysql.TypeLongBlob,
			pg.TypeCidr:               mysql.TypeVarchar, // 43
			pg.TypeInet:               mysql.TypeVarchar, // 43
			pg.TypeMacaddr:            mysql.TypeVarchar, // 17
			pg.TypeUUID:               mysql.TypeVarchar, // 36
			pg.TypeXML:                mysql.TypeLongText,
			pg.TypeJSON:               mysql.TypeLongText,
			pg.TypeTSVector:           mysql.TypeLongText,
			pg.TypeTSQuery:            mysql.TypeLongText,
			pg.TypeArray:              mysql.TypeLongText,
			pg.TypePoint:              mysql.TypeLongText,
			pg.TypeLine:               mysql.TypeLineString,
			pg.TypeLseq:               mysql.TypeLineString,
			pg.TypeBox:                mysql.TypePolygon,
			pg.TypePath:               mysql.TypeLineString,
			pg.TypePolygon:            mysql.TypePolygon,
			pg.TypeCircle:             mysql.TypePolygon,
			pg.TxidSnapshot:           mysql.TypeVarchar,
		},
	},
}
