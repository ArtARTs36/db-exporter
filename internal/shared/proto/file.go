package proto

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/iox"
	"github.com/artarts36/gds"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"log/slog"
)

type File struct {
	Package  string
	Services []*Service
	Messages []*Message
	Enums    []*Enum
	Imports  *gds.Set[string]
	Options  map[string]Option
}

type Option struct {
	Value  string
	Quotes bool
}

func PrepareOptions(options orderedmap.OrderedMap[string, interface{}]) map[string]Option {
	opts := map[string]Option{}

	for pair := options.Oldest(); pair != nil; pair = pair.Next() {
		key, val := pair.Key, pair.Value

		opt := Option{}

		switch v := val.(type) {
		case string:
			opt.Value = v

			if key != "optimize_for" {
				opt.Quotes = true
			}
		case bool:
			if v {
				opt.Value = "true"
			} else {
				opt.Value = "false"
			}
		case int, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float32, float64:
			opt.Value = fmt.Sprintf("%d", v)
		default:
			slog.Warn("[proto][prepare-options] unable prepare option", slog.String("key", key))
		}

		opts[key] = opt
	}

	return opts
}

func (f *File) Render(buf iox.Writer) {
	f.writeSyntax(buf)
	f.writePackage(buf)
	f.writeImports(buf)
	f.writeOptions(buf)
	f.writeServices(buf)
	f.writeMessages(buf)
	f.writeEnums(buf)
}

func (f *File) writeSyntax(buf iox.Writer) {
	buf.WriteString("syntax = \"proto3\";\n")
}

func (f *File) writeImports(buf iox.Writer) {
	if f.Imports == nil || f.Imports.Len() == 0 {
		return
	}

	buf.WriteString("\n")

	for _, im := range f.Imports.List() {
		buf.WriteString("import \"" + im + "\";\n")
	}
}

func (f *File) writeOptions(buf iox.Writer) {
	if len(f.Options) > 0 {
		buf.WriteString("\n")
	}

	for optName, opt := range f.Options {
		if opt.Quotes {
			buf.WriteString("option " + optName + " = \"" + opt.Value + "\";\n")
		} else {
			buf.WriteString("option " + optName + " = " + opt.Value + ";\n")
		}
	}
}

func (f *File) writeServices(buf iox.Writer) {
	for _, service := range f.Services {
		buf.WriteNewLine()
		service.write(buf)
	}
}

func (f *File) writePackage(buf iox.Writer) {
	if f.Package == "" {
		return
	}

	buf.WriteString("\npackage " + f.Package + ";\n")
}

func (f *File) writeMessages(buf iox.Writer) {
	for i, message := range f.Messages {
		buf.WriteNewLine()
		message.write(buf)

		if i < len(f.Messages)-1 {
			buf.WriteString("\n")
		}
	}
}

func (f *File) writeEnums(buf iox.Writer) {
	if len(f.Enums) == 0 {
		return
	}

	buf.WriteString("\n")

	for _, enum := range f.Enums {
		buf.WriteNewLine()
		enum.write(buf)
	}
}
