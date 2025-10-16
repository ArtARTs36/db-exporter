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
	}}
}

func (f *File) AddImport(dependency string) {
	f.Imports.Add(dependency)
}

func (f *File) SetOptions(options map[string]proto.Option) *File {
	f.Options = options
	return f
}
