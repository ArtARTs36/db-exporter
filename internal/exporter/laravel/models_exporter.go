package laravel

import (
	"context"
	"errors"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/php"
	"github.com/tyler-sommer/stick"
)

type ModelsExporter struct {
	pager *common.Pager
}

type laravelModel struct {
	Name       string
	Table      string
	Properties []*laravelModelProperty
	Dates      []string
	PrimaryKey laravelModelPrimaryKey
}

type laravelModelPrimaryKey struct {
	Exists bool

	Name string

	IsMultiple bool

	Column       string
	Type         string
	Incrementing bool
}

type laravelModelProperty struct {
	Name string
	Type string
}

type laravelModelSchema struct {
	Models []*laravelModel
}

func NewLaravelModelsExporter(pager *common.Pager) *ModelsExporter {
	return &ModelsExporter{
		pager: pager,
	}
}

func (e *ModelsExporter) ExportPerFile(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.LaravelModelsExportSpec)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())
	namespace := e.selectNamespace(spec)

	modelPage := e.pager.Of("laravel/model.php")

	for _, table := range params.Schema.Tables.List() {
		laravelSch := e.makeLaravelModelSchema([]*schema.Table{
			table,
		}, spec)

		page, err := modelPage.Export(
			fmt.Sprintf("%s.php", table.Name.Singular().Pascal()),
			map[string]stick.Value{
				"schema":    laravelSch,
				"namespace": namespace,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to render: %w", err)
		}

		pages = append(pages, page)
	}

	return pages, nil
}

func (e *ModelsExporter) Export(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.LaravelModelsExportSpec)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	namespace := e.selectNamespace(spec)

	laravelSch := e.makeLaravelModelSchema(params.Schema.Tables.List(), spec)

	page, err := e.pager.Of("laravel/model.php").Export(
		"models.php",
		map[string]stick.Value{
			"schema":    laravelSch,
			"namespace": namespace,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to render: %w", err)
	}

	return []*exporter.ExportedPage{
		page,
	}, nil
}

func (e *ModelsExporter) selectNamespace(spec *config.LaravelModelsExportSpec) string {
	if spec.Namespace != "" {
		return spec.Namespace
	}

	return "App\\Models"
}

func (e *ModelsExporter) makeLaravelModelSchema(
	tables []*schema.Table,
	spec *config.LaravelModelsExportSpec,
) *laravelModelSchema {
	modelSchema := &laravelModelSchema{
		Models: make([]*laravelModel, len(tables)),
	}

	for i, table := range tables {
		model := &laravelModel{
			Name:       table.Name.Singular().Pascal().Value,
			Table:      table.Name.Value,
			Properties: make([]*laravelModelProperty, 0, len(table.Columns)),
			Dates:      []string{},
			PrimaryKey: e.createModelPrimaryKey(table),
		}

		for _, column := range table.Columns {
			model.Properties = append(model.Properties, &laravelModelProperty{
				Name: column.Name.Value,
				Type: e.mapPhpType(column, model, spec),
			})
		}

		modelSchema.Models[i] = model
	}

	return modelSchema
}

func (*ModelsExporter) createModelPrimaryKey(table *schema.Table) laravelModelPrimaryKey {
	pk := table.PrimaryKey

	if pk == nil {
		return laravelModelPrimaryKey{
			Exists: false,
		}
	}

	if pk.ColumnsNames.Len() > 1 {
		return laravelModelPrimaryKey{
			Exists:     true,
			Name:       pk.Name.Value,
			IsMultiple: true,
		}
	}

	mapType := func(col *schema.Column) php.Type {
		switch col.PreparedType { //nolint:exhaustive // not need
		case schema.ColumnTypeInteger, schema.ColumnTypeInteger16, schema.ColumnTypeInteger64:
			return php.TypeInt
		case schema.ColumnTypeFloat32, schema.ColumnTypeFloat64:
			return php.TypeFloat
		case schema.ColumnTypeString, schema.ColumnTypeBytes:
			return php.TypeString
		}

		return php.TypeUndefined
	}

	pkColumnName := pk.ColumnsNames.First()

	var pkColumn *schema.Column

	for _, column := range table.Columns {
		if column.Name.Value == pkColumnName {
			pkColumn = column
			break
		}
	}

	if pkColumn == nil {
		return laravelModelPrimaryKey{Exists: false}
	}

	pkColType := mapType(pkColumn)
	if pkColType == php.TypeUndefined {
		return laravelModelPrimaryKey{Exists: false}
	}

	lpk := laravelModelPrimaryKey{
		Exists:       true,
		Name:         pk.Name.Value,
		IsMultiple:   false,
		Column:       pkColumnName,
		Type:         pkColType.String(),
		Incrementing: pkColType == php.TypeInt, // @todo need another way
	}

	return lpk
}

func (*ModelsExporter) mapPhpType(
	col *schema.Column,
	model *laravelModel,
	spec *config.LaravelModelsExportSpec,
) string {
	switch col.PreparedType {
	case schema.ColumnTypeInteger, schema.ColumnTypeInteger16, schema.ColumnTypeInteger64:
		return php.TypeInt.String()
	case schema.ColumnTypeFloat32, schema.ColumnTypeFloat64:
		return php.TypeFloat.String()
	case schema.ColumnTypeString, schema.ColumnTypeBytes:
		return php.TypeString.String()
	case schema.ColumnTypeBoolean:
		return php.TypeBool.String()
	case schema.ColumnTypeTimestamp:
		if !col.Name.Equal("created_at") && !col.Name.Equal("updated_at") {
			model.Dates = append(model.Dates, col.Name.Value)
		}

		if spec.TimeAs == "datetime" {
			return `\DateTimeInterface`
		}

		return `\Illuminate\Support\Carbon`
	}

	return ""
}