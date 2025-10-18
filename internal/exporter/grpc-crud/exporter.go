package grpccrud

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/modifiers"
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/indentx"
	"github.com/artarts36/db-exporter/internal/shared/proto"
)

type Exporter struct{}

type buildProcedureContext struct {
	sourceDriver config.DatabaseDriver

	service *presentation.Service

	tableSingularName string
	enumPages         map[string]*exporter.ExportedPage
}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) ExportPerFile(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.GRPCCrudExportSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec")
	}

	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len()+len(params.Schema.Enums))
	options := proto.PrepareOptions(spec.Options)
	enumPages := map[string]*exporter.ExportedPage{}
	indent := indentx.NewIndent(spec.Indent)

	pkg := e.newPackage(spec)

	for _, enum := range params.Schema.Enums {
		if enum.Used == 0 || enum.UsingInSingleTable() {
			continue
		}

		prfile := pkg.CreateFile(fmt.Sprintf("%s_enum.proto", enum.Name.Value)).SetOptions(options)
		prfile.AddEnum(*enum.Name, enum.Values)

		expPage := &exporter.ExportedPage{
			FileName: prfile.Name(),
			Content:  []byte(prfile.Render(indent)),
		}

		enumPages[enum.Name.Value] = expPage
		pages = append(pages, expPage)
	}

	for _, table := range params.Schema.Tables.List() {
		prfile := pkg.CreateFile(fmt.Sprintf("%s.proto", table.Name.Snake().Lower())).SetOptions(options)

		err := e.buildService(params.Schema.Driver, prfile, table, enumPages)
		if err != nil {
			return nil, fmt.Errorf("build service for table %q: %w", table.Name, err)
		}

		for _, enum := range table.UsingEnums {
			if !enum.UsingInSingleTable() {
				continue
			}

			prfile.AddEnum(*enum.Name, enum.Values)
		}

		expPage := &exporter.ExportedPage{
			FileName: prfile.Name(),
			Content:  []byte(prfile.Render(indent)),
		}

		pages = append(pages, expPage)
	}

	return pages, nil
}

func (e *Exporter) newPackage(spec *config.GRPCCrudExportSpec) *presentation.Package {
	configurators := []presentation.Configurator{}

	if spec.With.Object != nil {
		if spec.With.Object.GoogleAPIFieldBehavior.Object != nil {
			fb := modifiers.GoogleAPIFieldBehavior{}

			configurators = append(configurators, presentation.WithModifyField(fb.ModifyField))
		}

		if spec.With.Object.GoogleApiHTTP.Object != nil {
			gh := modifiers.GoogleApiHttp{
				PathPrefix: spec.With.Object.GoogleApiHTTP.Object.PathPrefix,
			}

			configurators = append(configurators, presentation.WithModifyProcedure(gh.ModifyProcedure))
		}
	}

	return presentation.NewPackage(spec.Package, configurators...)
}

func (e *Exporter) Export(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.GRPCCrudExportSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec, expected GRPCCrudExportSpec, got %T", params.Spec)
	}

	options := proto.PrepareOptions(spec.Options)

	pkg := e.newPackage(spec)

	prfile := pkg.CreateFile("services.proto").SetOptions(options)

	for _, enum := range params.Schema.Enums {
		if enum.Used == 0 {
			continue
		}

		prfile.AddEnum(*enum.Name, enum.Values)
	}

	for _, table := range params.Schema.Tables.List() {
		err := e.buildService(params.Schema.Driver, prfile, table, map[string]*exporter.ExportedPage{})
		if err != nil {
			return nil, fmt.Errorf("build service for table %q: %w", table.Name.Value, err)
		}
	}

	expPage := &exporter.ExportedPage{
		FileName: "services.proto",
		Content:  []byte(prfile.Render(indentx.NewIndent(spec.Indent))),
	}

	return []*exporter.ExportedPage{
		expPage,
	}, nil
}

func (e *Exporter) buildService(
	sourceDriver config.DatabaseDriver,
	prfile *presentation.File,
	table *schema.Table,
	enumPages map[string]*exporter.ExportedPage,
) error {
	procedureBuilders := map[presentation.ProcedureType]func(buildCtx *buildProcedureContext) error{
		presentation.ProcedureTypeList:   e.buildListProcedure,
		presentation.ProcedureTypeGet:    e.buildGetProcedure,
		presentation.ProcedureTypeDelete: e.buildDeleteProcedure,
		presentation.ProcedureTypeCreate: e.buildCreateProcedure,
		presentation.ProcedureTypePatch:  e.buildPatchProcedure,
	}

	mapColumnType := func(col *schema.Column) string {
		return e.mapType(sourceDriver, col, prfile, enumPages)
	}

	srv := prfile.AddService(
		table,
		func(message *presentation.TableMessage) {
			for _, column := range table.Columns {
				creator := func(field *presentation.Field) {
					field.SetType(mapColumnType(column))

					if !column.Nullable {
						field.AsRequired()
					}
				}

				fieldName := column.Name.Snake().Lower().Value

				if column.IsPrimaryKey() {
					message.CreatePrimaryKeyField(fieldName, column.Name.Value, creator)
				} else {
					message.CreateField(fieldName, column.Name.Value, creator)
				}
			}
		},
	)

	buildCtx := &buildProcedureContext{
		sourceDriver: sourceDriver,
		service:      srv,
		enumPages:    enumPages,
	}
	buildCtx.tableSingularName = buildCtx.service.TableMessage().Name()

	for procType, builder := range procedureBuilders {
		err := builder(buildCtx)
		if err != nil {
			return fmt.Errorf("build procedure of %s: %w", string(procType), err)
		}
	}

	return nil
}

func (e *Exporter) buildGetProcedure(
	buildCtx *buildProcedureContext,
) error {
	if buildCtx.service.TableMessage().PrimaryKey == nil {
		return nil
	}

	buildCtx.service.AddProcedureFn(
		"Get",
		presentation.ProcedureTypeGet,
		func(message *presentation.Message) {
			message.SetName(fmt.Sprintf("Get%sRequest", buildCtx.tableSingularName))

			for _, pkField := range buildCtx.service.TableMessage().PrimaryKey {
				message.CreateField(pkField.Name(), func(field *presentation.Field) {
					field.CopyType(field).AsRequired()
				})
			}
		},
		func(message *presentation.Message) {
			message.
				SetName(fmt.Sprintf("Get%sResponse", buildCtx.tableSingularName)).
				CreateField(
					buildCtx.tableSingularName,
					func(field *presentation.Field) {
						field.SetType(buildCtx.service.TableMessage().Name()).AsRequired()
					},
				)
		},
	)

	return nil
}

func (e *Exporter) buildListProcedure(
	buildCtx *buildProcedureContext,
) error {
	buildCtx.service.AddProcedureFn("List", presentation.ProcedureTypeList,
		func(message *presentation.Message) {
			message.SetName(fmt.Sprintf("List%sRequest", buildCtx.tableSingularName))
		},
		func(message *presentation.Message) {
			message.
				SetName(fmt.Sprintf("List%sResponse", buildCtx.tableSingularName)).
				CreateField("items", func(field *presentation.Field) {
					field.AsRepeated().SetType(buildCtx.service.TableMessage().Name())
				})
		},
	)

	return nil
}

func (e *Exporter) buildDeleteProcedure(
	buildCtx *buildProcedureContext,
) error {
	if buildCtx.service.TableMessage().Table.PrimaryKey == nil {
		return nil
	}

	var err error

	buildCtx.service.AddProcedureFn(
		"Delete",
		presentation.ProcedureTypeDelete,
		func(message *presentation.Message) {
			message.SetName(fmt.Sprintf("Delete%sRequest", buildCtx.tableSingularName))

			for _, pkField := range buildCtx.service.TableMessage().PrimaryKey {
				message.CreateField(pkField.Name(), func(field *presentation.Field) {
					field.CopyType(pkField)
					field.AsRequired()
				})
			}
		},
		func(message *presentation.Message) {
			message.SetName(fmt.Sprintf("Delete%sResponse", buildCtx.tableSingularName))
		},
	)

	return err
}

func (e *Exporter) buildCreateProcedure(
	buildCtx *buildProcedureContext,
) error {
	if buildCtx.service.TableMessage().Table.PrimaryKey == nil {
		return nil
	}

	buildCtx.service.AddProcedureFn("Create", presentation.ProcedureTypeCreate,
		func(message *presentation.Message) {
			message.SetName(fmt.Sprintf("Create%sRequest", buildCtx.tableSingularName))

			for _, col := range buildCtx.service.TableMessage().Table.Columns {
				if col.IsAutoincrement {
					continue
				}

				tableField, _ := message.Service().TableMessage().GetField(col.Name.Value)

				message.CreateField(tableField.Name(), func(field *presentation.Field) {
					field.CopyType(tableField)

					if !col.Nullable {
						field.AsRequired()
					}
				})
			}
		},
		func(message *presentation.Message) {
			message.SetName(fmt.Sprintf("Create%sResponse", buildCtx.tableSingularName))
		},
	)

	return nil
}

func (e *Exporter) buildPatchProcedure(
	buildCtx *buildProcedureContext,
) error {
	buildCtx.service.AddProcedureFn("Patch", presentation.ProcedureTypePatch,
		func(message *presentation.Message) {
			message.SetName(fmt.Sprintf("Patch%sRequest", buildCtx.tableSingularName))

			for _, col := range buildCtx.service.TableMessage().Table.Columns {
				tableField, _ := message.Service().TableMessage().GetField(col.Name.Value)

				message.CreateField(tableField.Name(), func(field *presentation.Field) {
					field.CopyType(tableField)

					if !col.Nullable {
						field.AsRequired()
					}
				})
			}
		},
		func(message *presentation.Message) {
			message.
				SetName(fmt.Sprintf("Patch%sResponse", buildCtx.tableSingularName)).
				CreateField(buildCtx.service.TableMessage().Table.Name.Pascal().Singular().Value, func(field *presentation.Field) {
					field.SetType(buildCtx.service.TableMessage().Name())
				})
		},
	)

	return nil
}
