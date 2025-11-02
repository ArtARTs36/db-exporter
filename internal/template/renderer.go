package template

import (
	"bytes"
	"github.com/tyler-sommer/stick"
	"github.com/tyler-sommer/stick/twig/filter"
	"io"
)

type Renderer struct {
	engine *stick.Env
}

func NewRenderer(templateLoader stick.Loader) *Renderer {
	eng := stick.New(templateLoader)
	eng.Functions = twigFuncs()
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

func (r *Renderer) RenderTo(name string, params map[string]stick.Value, w io.Writer) error {
	r.extendParams(params)

	err := r.engine.Execute(name, w, params)
	if err != nil {
		return err
	}

	return err
}

func (r *Renderer) extendParams(params map[string]stick.Value) {
	params["figure_brace_opened"] = "{"
}
