package presentation

import (
	"github.com/artarts36/db-exporter/internal/shared/proto"
	"github.com/artarts36/gds"
)

type File struct {
	proto.File
}

func NewFile(pkg string) *File {
	return &File{proto.File{
		Package:  pkg,
		Imports:  gds.NewSet[string](),
		Services: make([]*proto.Service, 0),
		Messages: make([]*proto.Message, 0),
		Enums:    make([]*proto.Enum, 0),
	}}
}

func AllocateFile(pkg string, enumsLength int) *File {
	return &File{proto.File{
		Package:  pkg,
		Imports:  gds.NewSet[string](),
		Services: make([]*proto.Service, 0),
		Messages: make([]*proto.Message, 0),
		Enums:    make([]*proto.Enum, 0, enumsLength),
	}}
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
