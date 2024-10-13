package template

import (
	"bytes"
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/ds"
	"strings"
	"time"

	"github.com/tyler-sommer/stick"
	"github.com/tyler-sommer/stick/twig/filter"
)

type Renderer struct {
	engine *stick.Env
}

func NewRenderer(templateLoader stick.Loader) *Renderer {
	eng := stick.New(templateLoader)
	eng.Functions["bool_string"] = func(_ stick.Context, args ...stick.Value) stick.Value {
		val, _ := args[0].(bool)
		if val {
			return "true"
		}
		return "false"
	}
	eng.Functions["spaces"] = func(_ stick.Context, args ...stick.Value) stick.Value {
		count, valid := args[0].(int)
		if !valid || count < 0 {
			count = 0
		}

		return strings.Repeat(" ", count)
	}
	eng.Functions["quote_string"] = func(_ stick.Context, args ...stick.Value) stick.Value {
		switch val := args[0].(type) {
		case string:
			return fmt.Sprintf("%q", val)
		case time.Time:
			if val.IsZero() {
				return ""
			}

			return fmt.Sprintf("%q", val.Format(time.RFC3339))
		}
		return args[0]
	}

	const spacesAfterArgsCount = 2
	eng.Functions["spaces_after"] = func(_ stick.Context, args ...stick.Value) stick.Value {
		if len(args) != spacesAfterArgsCount {
			return ""
		}

		var currentStringLen int

		switch v := args[0].(type) {
		case string:
			currentStringLen = len(v)
		case *ds.String:
			currentStringLen = v.Len()
		default:
			return ""
		}

		needLength, valid := args[1].(int)
		if !valid {
			return ""
		}

		repeats := needLength - currentStringLen
		if repeats == 0 {
			return ""
		}

		return strings.Repeat(" ", repeats)
	}

	eng.Filters = filter.TwigFilters()

	return &Renderer{
		engine: eng,
	}
}

func (r *Renderer) Render(name string, params map[string]stick.Value) ([]byte, error) {
	buf := &bytes.Buffer{}

	r.extendParams(params)

	err := r.engine.Execute(name, buf, params)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), err
}

func (r *Renderer) extendParams(params map[string]stick.Value) {
	params["figure_brace_opened"] = "{"
}
