package proto

import (
	"fmt"
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
