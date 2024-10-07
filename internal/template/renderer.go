package template

import (
	"bytes"
	"fmt"
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
		val := args[0].(bool)
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

	eng.Filters = filter.TwigFilters()

	return &Renderer{
		engine: eng,
	}
}

func (r *Renderer) Render(name string, params map[string]stick.Value) ([]byte, error) {
	buf := &bytes.Buffer{}

	err := r.engine.Execute(name, buf, params)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), err
}
