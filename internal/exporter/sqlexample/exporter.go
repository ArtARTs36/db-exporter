package sqlexample

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/schema"
)

type Exporter struct {
	valueResolvers []valueResolver
}

func NewExporter() *Exporter {
	return &Exporter{
		valueResolvers: make([]valueResolver, 0),
	}
}

func newExporter(valueResolvers []valueResolver) *Exporter {
	return &Exporter{valueResolvers: valueResolvers}
}

func (e *Exporter) ExportPerFile(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {

}

func (e *Exporter) Export(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	return nil, nil
}

type example struct {
	table *schema.Table

	Values map[string]interface{}
}

func (e *Exporter) genExamples(table *schema.Table, count int) ([]*example, error) {
	examples := make([]*example, count)

	resolvers := map[string]valueResolver{}

	for _, column := range table.Columns {
		resolverFound := false

		for _, resolver := range e.valueResolvers {
			if resolver.supports(column) {
				resolvers[column.Name.Value] = resolver
				resolverFound = true
				continue
			}
		}

		if !resolverFound {
			return nil, fmt.Errorf("no value resolver found for column %q.%q", table.Name.Value, column.Name)
		}
	}

}
