package grpccrud

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/modifiers"
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/paginator"
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"
	"github.com/artarts36/db-exporter/internal/infrastructure/workspace"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/iox"
	"github.com/artarts36/db-exporter/internal/shared/proto"
	"github.com/artarts36/gds"
)

type Exporter struct{}

type buildProcedureContext struct {
	sourceDriver config.DatabaseDriver

	service *presentation.Service

	paginator paginator.Paginator

	tableSingularName string
	tablePluralName   string
}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) createPaginator(spec *config.GRPCCrudExportSpec) paginator.Paginator {
	if spec.Pagination == config.PaginationTypeToken {
		return &paginator.Token{}
	}

	if spec.Pagination == config.PaginationTypeNone {
		return &paginator.None{}
	}

	return &paginator.Offset{}
}

func (e *Exporter) ExportPerFile(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.GRPCCrudExportSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec")
	}

	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len()+len(params.Schema.Enums))
	options := proto.PrepareOptions(spec.Options)
	indent := iox.NewIndent(spec.Indent)

	pkg := e.newPackage(spec)

	for _, enum := range params.Schema.Enums {
		if enum.Used == 0 || enum.UsingInSingleTable() {
			continue
		}

		err := params.Workspace.Write(
			ctx,
			fmt.Sprintf("%s_enum.proto", enum.Name.Value),
			func(buffer workspace.Buffer,
			) error {
				prfile := pkg.CreateFile(fmt.Sprintf("%s_enum.proto", enum.Name.Value)).SetOptions(options)
				prfile.AddEnum(*enum.Name, enum.Values)
				prfile.Render(buffer, indent)

				return nil
			})
		if err != nil {
			return nil, fmt.Errorf("write enum to workspace: %w", err)
		}
	}

	pager := e.createPaginator(spec)

	for _, table := range params.Schema.Tables.List() {
		prfile := pkg.CreateFile(fmt.Sprintf("%s.proto", table.Name.Snake().Lower())).SetOptions(options)

		err := params.Workspace.Write(ctx, prfile.Name(), func(buffer workspace.Buffer) error {
			err := e.buildService(params.Schema.Driver, prfile, table, pager)
			if err != nil {
				return fmt.Errorf("build service for table %q: %w", table.Name, err)
			}

			for _, enum := range table.UsingEnums {
				if !enum.UsingInSingleTable() {
					continue
				}

				prfile.AddEnum(*enum.Name, enum.Values)
			}

			prfile.Render(buffer, indent)

			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("write table %q to workspace: %w", table.Name.Value, err)
		}
	}

	return pages, nil
}

func (e *Exporter) newPackage(spec *config.GRPCCrudExportSpec) *presentation.Package {
	configurators := []presentation.Configurator{}

	commentsModifier := &modifiers.Comments{}

	configurators = append(
		configurators,
		presentation.WithModifyProcedure(commentsModifier.ModifyProcedure),
		presentation.WithModifyService(commentsModifier.ModifyService),
	)

	if spec.With.Object != nil {
		if spec.With.Object.GoogleAPIFieldBehavior.Object != nil {
			fb := modifiers.GoogleAPIFieldBehavior{}

			configurators = append(configurators, presentation.WithModifyField(fb.ModifyField))
		}

		if spec.With.Object.GoogleApiHttp.Object != nil {
			gh := modifiers.GoogleApiHttp{
				PathPrefix: spec.With.Object.GoogleApiHttp.Object.PathPrefix,
			}

			configurators = append(configurators, presentation.WithModifyProcedure(gh.ModifyProcedure))
		}

		if spec.With.Object.BufValidateField.Object != nil {
			bufValidateField := modifiers.BufValidate{}

			configurators = append(configurators, presentation.WithModifyField(bufValidateField.ModifyField))
		}
	}

	return presentation.NewPackage(spec.Package, configurators...)
}

func (e *Exporter) Export(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.GRPCCrudExportSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec, expected GRPCCrudExportSpec, got %T", params.Spec)
	}

	options := proto.PrepareOptions(spec.Options)

	pkg := e.newPackage(spec)
	pager := e.createPaginator(spec)
	indent := iox.NewIndent(spec.Indent)

	prfile := pkg.CreateFile("services.proto").SetOptions(options)

	err := params.Workspace.Write(ctx, "services.proto", func(buffer workspace.Buffer) error {
		for _, enum := range params.Schema.Enums {
			if enum.Used == 0 {
				continue
			}

			prfile.AddEnum(*enum.Name, enum.Values)
		}

		for _, table := range params.Schema.Tables.List() {
			err := e.buildService(params.Schema.Driver, prfile, table, pager)
			if err != nil {
				return fmt.Errorf("build service for table %q: %w", table.Name.Value, err)
			}
		}

		expPage := &exporter.ExportedPage{
			FileName: "services.proto",
			Content:  []byte(prfile.Render(iox.NewWriterWithIndent(spec.Indent))),
		}

		return nil
	})

	return nil, err
}

func (e *Exporter) buildService(
	sourceDriver config.DatabaseDriver,
	prfile *presentation.File,
	table *schema.Table,
	pager paginator.Paginator,
) error {
	procedureBuilders := []struct {
		Type  presentation.ProcedureType
		Build func(buildCtx *buildProcedureContext) error
	}{
		{
			Type:  presentation.ProcedureTypeList,
			Build: e.buildListProcedure,
		},
		{
			Type:  presentation.ProcedureTypeGet,
			Build: e.buildGetProcedure,
		},
		{
			Type:  presentation.ProcedureTypeDelete,
			Build: e.buildDeleteProcedure,
		},
		{
			Type:  presentation.ProcedureTypeCreate,
			Build: e.buildCreateProcedure,
		},
		{
			Type:  presentation.ProcedureTypePatch,
			Build: e.buildPatchProcedure,
		},
	}

	mapColumnType := func(col *schema.Column) string {
		return e.mapType(sourceDriver, col, prfile)
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

					field.SetColumn(column)
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
		paginator:    pager,
	}
	buildCtx.tableSingularName = buildCtx.service.TableMessage().Name()
	buildCtx.tablePluralName = buildCtx.service.TableMessage().Table().Name.Pascal().Plural().Value

	for _, builder := range procedureBuilders {
		err := builder.Build(buildCtx)
		if err != nil {
			return fmt.Errorf("build procedure of %s: %w", string(builder.Type), err)
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
					field.CopyType(pkField).AsRequired()
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
			message.SetName(fmt.Sprintf("List%sRequest", buildCtx.tablePluralName))

			if len(message.Service().TableMessage().PrimaryKey) == 1 {
				message.CreateField(
					gds.NewString(message.Service().TableMessage().PrimaryKey[0].Name()).Plural().Value,
					func(field *presentation.Field) {
						field.
							CopyType(message.Service().TableMessage().PrimaryKey[0]).
							AsRepeated().
							NotRequired()
					},
				)
			}

			buildCtx.paginator.AddPaginationToRequest(message)
		},
		func(message *presentation.Message) {
			message.
				SetName(fmt.Sprintf("List%sResponse", buildCtx.tablePluralName)).
				CreateField("items", func(field *presentation.Field) {
					field.AsRepeated().SetType(buildCtx.service.TableMessage().Name())
				})

			buildCtx.paginator.AddPaginationToResponse(message)
		},
	)

	return nil
}

func (e *Exporter) buildDeleteProcedure(
	buildCtx *buildProcedureContext,
) error {
	if buildCtx.service.TableMessage().Table().PrimaryKey == nil {
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
	if buildCtx.service.TableMessage().Table().PrimaryKey == nil {
		return nil
	}

	buildCtx.service.AddProcedureFn("Create", presentation.ProcedureTypeCreate,
		func(message *presentation.Message) {
			message.SetName(fmt.Sprintf("Create%sRequest", buildCtx.tableSingularName))

			for _, col := range buildCtx.service.TableMessage().Table().Columns {
				if e.columnAutofilled(col) {
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
			message.
				SetName(fmt.Sprintf("Create%sResponse", buildCtx.tableSingularName)).
				CreateField(
					buildCtx.service.TableMessage().Table().Name.Pascal().Singular().Value,
					func(field *presentation.Field) {
						field.SetType(buildCtx.service.TableMessage().Name()).AsRequired()
					},
				)
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

			for _, col := range buildCtx.service.TableMessage().Table().Columns {
				if e.columnAutofilled(col) {
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
			message.
				SetName(fmt.Sprintf("Patch%sResponse", buildCtx.tableSingularName)).
				CreateField(
					buildCtx.service.TableMessage().Table().Name.Pascal().Singular().Value,
					func(field *presentation.Field) {
						field.SetType(buildCtx.service.TableMessage().Name()).AsRequired()
					},
				)
		},
	)

	return nil
}

func (e *Exporter) columnAutofilled(col *schema.Column) bool {
	if col.IsAutoincrement {
		return true
	}

	if !col.DefaultRaw.Valid {
		return false
	}

	return col.Name.Equal("id", "created_at", "updated_at", "deleted_at")
}
