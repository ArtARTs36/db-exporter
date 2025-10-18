package grpccrud

import (
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"
	"github.com/artarts36/db-exporter/internal/infrastructure/sqltype"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/golang"
)

func (e *Exporter) mapType(
	sourceDriver config.DatabaseDriver,
	column *schema.Column,
	file *presentation.File,
	enumPages map[string]*exporter.ExportedPage,
) string {
	if column.Enum != nil {
		enumPage, enumPageExists := enumPages[column.Enum.Name.Value]
		if enumPageExists {
			file.AddImport(enumPage.FileName)
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
	default:
		return "string"
	}
}
