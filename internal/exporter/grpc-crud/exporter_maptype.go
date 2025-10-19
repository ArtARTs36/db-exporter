package grpccrud

import (
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"
	"github.com/artarts36/db-exporter/internal/infrastructure/sqltype"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/golang"
)

func (e *Exporter) mapType(
	sourceDriver config.DatabaseDriver,
	column *schema.Column,
	file *presentation.File,
) string {
	if column.Enum != nil {
		enumFilename, enumPageExists := file.Package().LocateEnum(column.Enum.Name.Value)
		if enumPageExists && enumFilename != file.Name() {
			file.AddImport(enumFilename)
		}

		return column.Enum.Name.Pascal().Value
	}

	goType := sqltype.MapGoType(sourceDriver, column.Type)

	switch goType {
	case golang.TypeInt, golang.TypeInt16, golang.TypeInt64:
		return "int64"
	case golang.TypeFloat64:
		return "double"
	case golang.TypeFloat32:
		return "double"
	case golang.TypeBool:
		return "bool"
	case golang.TypeTimeTime:
		file.AddImport("google/protobuf/timestamp.proto")

		return "google.protobuf.Timestamp"
	case golang.TypeTimeDuration:
		file.AddImport("google/protobuf/duration.proto")

		return "google.protobuf.Duration"
	default:
		return "string"
	}
}
