package bufvalidate

import "github.com/artarts36/db-exporter/internal/shared/proto"

func MaxLen(length string) *proto.FieldOption {
	return &proto.FieldOption{
		Name:  "(buf.validate.field).string.max_len",
		Value: proto.ConstValue(length),
	}
}

func UUID() *proto.FieldOption {
	return &proto.FieldOption{
		Name:  "(buf.validate.field).string.uuid",
		Value: true,
	}
}

func Email() *proto.FieldOption {
	return &proto.FieldOption{
		Name:  "(buf.validate.field).string.email",
		Value: true,
	}
}

func Required() *proto.FieldOption {
	return &proto.FieldOption{
		Name:  "(buf.validate.field).required",
		Value: true,
	}
}
