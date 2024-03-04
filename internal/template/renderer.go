package template

import (
	"bytes"
	"github.com/tyler-sommer/stick"
)

type Renderer struct {
	engine *stick.Env
}

func InitRenderer(rootDir string) *Renderer {
	return &Renderer{
		engine: stick.New(stick.NewFilesystemLoader(rootDir)),
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
