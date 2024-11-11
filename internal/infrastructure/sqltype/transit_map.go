package sqltype

import (
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/schema"
)

var transitSQLTypeMap = map[config.DatabaseDriver]map[config.DatabaseDriver]map[schema.Type]schema.Type{ //nolint:exhaustive,lll // not need
	config.DatabaseDriverDBML: {
		config.DatabaseDriverPostgres: {
			DBMLChar:    PGCharacter,
			DBMLVarchar: PGCharacterVarying,
			DBMLBinary:  PGBytea,
			DBMLText:    PGText,

			DBMLInt:     PGInteger,
			DBMLInteger: PGInteger,

			DBMLTimestamp: PGTimestampWithTZ,

			DBMLUUID: PGUUID,

			DBMLFloat:  PGDoublePrecision,
			DBMLFloat8: PGFloat8,
		},
	},
	config.DatabaseDriverPostgres: {
		config.DatabaseDriverDBML: {
			PGCharacter:        DBMLChar,
			PGCharacterVarying: DBMLVarchar,
			PGBytea:            DBMLBinary,
			PGText:             DBMLText,

			PGInteger: DBMLInteger,

			PGTimestampWithTZ:    DBMLTimestamp,
			PGTimestampWithoutTZ: DBMLTimestamp,

			PGUUID: DBMLUUID,

			PGDoublePrecision: DBMLFloat,
			PGFloat8:          PGFloat8,
		},

		// https://dev.mysql.com/doc/workbench/en/wb-migration-database-postgresql-typemapping.html
		config.DatabaseDriverMySQL: {
			PGInt:                MySQLInt,
			PGSmallInt:           MySQLSmallInt,
			PGBigint:             MySQLInt,
			PGSerial:             MySQLInt,
			PGSmallSerial:        MySQLSmallInt,
			PGBigSerial:          MySQLBigInt,
			PGBit:                MySQLBit,
			PGBoolean:            MySQLTinyint.WithLength("1"),
			PGReal:               MySQLFloat,
			PGDoublePrecision:    MySQLDouble,
			PGNumeric:            MySQLDecimal,
			PGDecimal:            MySQLDecimal,
			PGMoney:              MySQLDecimal.WithLength("19,2"),
			PGCharacter:          MySQLChar,
			PGChar:               MySQLChar,
			PGCharacterVarying:   MySQLLongText,
			PGDate:               MySQLDate,
			PGTimeWithTZ:         MySQLTime,
			PGTimeWithoutTZ:      MySQLTime,
			PGTimestampWithTZ:    MySQLDateTime,
			PGTimestampWithoutTZ: MySQLDateTime,
			PGInterval:           MySQLTime,
			PGBytea:              MySQLLongBlob,
			PGCidr:               MySQLVarchar.WithLength("43"), // 43
			PGInet:               MySQLVarchar.WithLength("43"), // 43
			PGMacaddr:            MySQLVarchar.WithLength("17"), // 17
			PGUUID:               MySQLVarchar.WithLength("36"), // 36
			PGXML:                MySQLLongText,
			PGJSON:               MySQLLongText,
			PGTSVector:           MySQLLongText,
			PGTSQuery:            MySQLLongText,
			PGArray:              MySQLLongText,
			PGPoint:              MySQLLongText,
			PGLine:               MySQLLineString,
			PGLseg:               MySQLLineString,
			PGBox:                MySQLPolygon,
			PGPath:               MySQLLineString,
			PGPolygon:            MySQLPolygon,
			PGCircle:             MySQLPolygon,
			PGTxidSnapshot:       MySQLVarchar,
		},
	},
}
