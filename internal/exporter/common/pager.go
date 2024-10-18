package common

import (
	"fmt"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/template"
)

type Pager struct {
	renderer *template.Renderer
}

type Page struct {
	renderer *template.Renderer
	tmplName string
}

func NewPager(renderer *template.Renderer) *Pager {
	return &Pager{
		renderer: renderer,
	}
}

func (p *Pager) Of(tmplName string) *Page {
	return &Page{
		renderer: p.renderer,
		tmplName: tmplName,
	}
}

func (p *Page) Export(filename string, params map[string]stick.Value) (*exporter.ExportedPage, error) {
	indexContent, err := p.renderer.Render(p.tmplName, params)
	if err != nil {
		return nil, fmt.Errorf("unable to render template %q: %w", p.tmplName, err)
	}

	return &exporter.ExportedPage{
		FileName: filename,
		Content:  indexContent,
	}, nil
}
