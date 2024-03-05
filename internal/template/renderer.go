package template

import (
	"bytes"
	"github.com/tyler-sommer/stick"
	"strings"
)

type Renderer struct {
	engine *stick.Env
}

func InitRenderer(rootDir string) *Renderer {
	eng := stick.New(stick.NewFilesystemLoader(rootDir))
	eng.Functions["spaces"] = func(ctx stick.Context, args ...stick.Value) stick.Value {
		count := args[0].(int)
		if count < 0 {
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
