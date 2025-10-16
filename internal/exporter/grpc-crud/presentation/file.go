package presentation

import (
	"github.com/artarts36/db-exporter/internal/shared/indentx"
	"github.com/artarts36/db-exporter/internal/shared/proto"
	"github.com/artarts36/gds"
)

type File struct {
	proto.File

	services []*Service
}

func NewFile(pkg string) *File {
	return &File{
		File: proto.File{
			Package:  pkg,
			Imports:  gds.NewSet[string](),
			Services: make([]*proto.Service, 0),
			Messages: make([]*proto.Message, 0),
			Enums:    make([]*proto.Enum, 0),
		},
		services: make([]*Service, 0), // @todo
	}
}

func AllocateFile(pkg string, enumsLength int) *File {
	return &File{
		File: proto.File{
			Package:  pkg,
			Imports:  gds.NewSet[string](),
			Services: make([]*proto.Service, 0),
			Messages: make([]*proto.Message, 0),
			Enums:    make([]*proto.Enum, 0, enumsLength),
		},
		services: make([]*Service, 0),
	}
}

func (f *File) AddImport(dependency string) {
	f.File.Imports.Add(dependency)
}

func (f *File) SetOptions(options map[string]proto.Option) *File {
	f.File.Options = options
	return f
}

func (f *File) AddEnum(name gds.String, values []string) {
	f.File.Enums = append(f.Enums, proto.NewEnumWithValues(name, values))
}

func (f *File) Render(indent *indentx.Indent) string {
	return f.File.Render(indent)
}

func (f *File) AddService(name string, procedures int) *Service {
	srv := &Service{
		Name:       name,
		Procedures: make([]*Procedure, 0, procedures),
	}

	return srv
}
