package presentation

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/indentx"
	"github.com/artarts36/db-exporter/internal/shared/proto"
	"github.com/artarts36/gds"
)

type File struct {
	name string
	pkg  *Package

	proto proto.File

	services []*Service

	cfg *config
}

func NewFile(pkg *Package, name string) *File {
	f := &File{
		name: name,
		pkg:  pkg,
		proto: proto.File{
			Package:  pkg.name,
			Imports:  gds.NewSet[string](),
			Services: make([]*proto.Service, 0),
			Messages: make([]*proto.Message, 0),
			Enums:    make([]*proto.Enum, 0),
		},
		services: make([]*Service, 0), // @todo
		cfg:      pkg.cfg,
	}

	return f
}

func (f *File) Name() string {
	return f.name
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

	f.pkg.registerEnumLocation(name.Value, f.name)
}

func (f *File) Render(indent *indentx.Indent) string {
	return f.proto.Render(indent)
}

func (f *File) AddService(
	table *schema.Table,
	createTableMessage func(*TableMessage),
) *Service {
	srv := &Service{
		proto: &proto.Service{
			Name:       fmt.Sprintf("%sService", table.Name.Pascal()),
			Procedures: make([]*proto.ServiceProcedure, 0),
		},
		file: f,
	}

	tableMsg := newTableMessage(table, srv)
	createTableMessage(tableMsg)

	srv.table = tableMsg

	f.AddMessage(tableMsg.message.proto)

	f.proto.Services = append(f.proto.Services, srv.proto)

	return srv
}

func (f *File) AddMessage(msg *proto.Message) {
	f.proto.Messages = append(f.proto.Messages, msg)
}

func (f *File) Package() *Package {
	return f.pkg
}
