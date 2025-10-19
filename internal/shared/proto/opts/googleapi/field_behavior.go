package googleapi

import "github.com/artarts36/db-exporter/internal/shared/proto"

func FieldRequired() *proto.FieldOption {
	return fieldBehavior("REQUIRED")
}

func FieldOutputOnly() *proto.FieldOption {
	return fieldBehavior("OUTPUT_ONLY")
}

func FieldOptional() *proto.FieldOption {
	return fieldBehavior("OPTIONAL")
}

func fieldBehavior(behavior string) *proto.FieldOption {
	return &proto.FieldOption{
		Name:  "(google.api.field_behavior)",
		Value: behavior,
	}
}
