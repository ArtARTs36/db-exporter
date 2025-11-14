package sqltype

import "github.com/artarts36/db-exporter/internal/schema"

func mapType(typeMap map[string]schema.DataType, name string) schema.DataType {
	t, ok := typeMap[name]
	if ok {
		return t
	}
	return schema.DataType{Name: name}
}
