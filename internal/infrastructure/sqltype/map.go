package sqltype

import "github.com/artarts36/db-exporter/internal/schema"

func mapType(typeMap map[string]schema.Type, name string) schema.Type {
	t, ok := typeMap[name]
	if ok {
		return t
	}
	return schema.Type{Name: name}
}
