package sqltype

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/schema"

	"github.com/artarts36/db-exporter/internal/config"
)

func TransitSQLType(source, target config.DatabaseDriver, dataType schema.Type) (schema.Type, error) {
	if source == target {
		return dataType, nil
	}

	sourceMap, ok := transitSQLTypeMap[source]
	if !ok {
		return schema.Type{}, fmt.Errorf("trasition map for source driver %q is not present", source)
	}

	targetMap, ok := sourceMap[target]
	if !ok {
		return schema.Type{}, fmt.Errorf("trasition map for target driver %q is not present", target)
	}

	t, ok := targetMap[dataType]
	if !ok {
		t = dataType
	}

	return t, nil
}
