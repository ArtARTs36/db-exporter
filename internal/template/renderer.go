package template

import (
	"bytes"
	"strings"

	"github.com/tyler-sommer/stick"
)

type Renderer struct {
	engine *stick.Env
}

func NewRenderer(templateLoader stick.Loader) *Renderer {
	eng := stick.New(templateLoader)
	eng.Functions["spaces"] = func(ctx stick.Context, args ...stick.Value) stick.Value {
		count, valid := args[0].(int)
		if !valid || count < 0 {
			count = 0
		}

		return strings.Repeat(" ", count)
	}

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
