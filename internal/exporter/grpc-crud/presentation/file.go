package presentation

import (
	"github.com/artarts36/db-exporter/internal/shared/indentx"
	"github.com/artarts36/db-exporter/internal/shared/proto"
	"github.com/artarts36/gds"
)

type File struct {
	proto proto.File

	services []*Service

	cfg *config
}

func NewFile(pkg string, configurators ...Configurator) *File {
	f := &File{
		proto: proto.File{
			Package:  pkg,
			Imports:  gds.NewSet[string](),
			Services: make([]*proto.Service, 0),
			Messages: make([]*proto.Message, 0),
			Enums:    make([]*proto.Enum, 0),
		},
		services: make([]*Service, 0), // @todo
		cfg:      newConfig(configurators),
	}

	return f
}

func AllocateFile(pkg string, enumsLength int, configurators ...Configurator) *File {
	return &File{
		proto: proto.File{
			Package:  pkg,
			Imports:  gds.NewSet[string](),
			Services: make([]*proto.Service, 0),
			Messages: make([]*proto.Message, 0),
			Enums:    make([]*proto.Enum, 0, enumsLength),
		},
		services: make([]*Service, 0),
		cfg:      newConfig(configurators),
	}
}

func (f *File) AddImport(dependency string) {
	f.proto.Imports.Add(dependency)
}

func (f *File) SetOptions(options map[string]proto.Option) *File {
	f.proto.Options = options
	return f
}

func (f *File) AddEnum(name gds.String, values []string) {
	f.proto.Enums = append(f.proto.Enums, proto.NewEnumWithValues(name, values))
}

func (f *File) Render(indent *indentx.Indent) string {
	return f.proto.Render(indent)
}

func (f *File) AddService(
	name string,
	tableMsg *TableMessage,
	procedures int,
) *Service {
	srv := &Service{
		table: tableMsg,
		file:  f,
		proto: &proto.Service{
			Name:       name,
			Procedures: make([]*proto.ServiceProcedure, 0, procedures),
		},
	}

	f.AddMessage(tableMsg.Proto)

	f.proto.Services = append(f.proto.Services, srv.proto)

	return srv
}

func (f *File) AddMessage(msg *proto.Message) {
	f.proto.Messages = append(f.proto.Messages, msg)
}
