package grpccrud

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/fieldmap"
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/modifiers"
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/tablemsg"
	"github.com/artarts36/db-exporter/internal/infrastructure/sqltype"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/golang"
	"github.com/artarts36/db-exporter/internal/shared/indentx"
	"github.com/artarts36/db-exporter/internal/shared/proto"
)

type Exporter struct{}

type buildProcedureContext struct {
	sourceDriver config.DatabaseDriver

	service *presentation.Service

	table             *schema.Table
	tableMsg          *presentation.TableMessage
	tableSingularName string
	enumPages         map[string]*exporter.ExportedPage
}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) newTableMapper(spec *config.GRPCCrudExportSpec) *tablemsg.Mapper {
	modifiers := []fieldmap.Modifier{}

	if spec.With.Object != nil {
		if spec.With.Object.GoogleAPIFieldBehavior.Object != nil {
			modifiers = append(modifiers, &fieldmap.GoogleAPIFieldBehaviorModifier{})
		}
	}

	if len(modifiers) == 0 {
		return tablemsg.NewMapper(fieldmap.Nop{})
	}

	if len(modifiers) == 1 {
		return tablemsg.NewMapper(modifiers[0])
	}

	return tablemsg.NewMapper(fieldmap.Compose(modifiers))
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

	for _, enum := range params.Schema.Enums {
		prfile := presentation.AllocateFile(spec.Package, 1).SetOptions(options)
		prfile.AddEnum(*enum.Name, enum.Values)

		expPage := &exporter.ExportedPage{
			FileName: fmt.Sprintf("%s_enum.proto", enum.Name.Value),
			Content:  []byte(prfile.Render(indent)),
		}

		enumPages[enum.Name.Value] = expPage
		pages = append(pages, expPage)
	}

	tablemsgMapper := e.newTableMapper(spec)

	for _, table := range params.Schema.Tables.List() {
		prfile := presentation.NewFile(spec.Package).SetOptions(options)

		srv, err := e.buildService(tablemsgMapper, params.Schema.Driver, prfile, table, enumPages)
		if err != nil {
			return nil, fmt.Errorf("build service for table %q: %w", table.Name, err)
		}
		if !srv.HasProcedures() {
			continue
		}

		expPage := &exporter.ExportedPage{
			FileName: fmt.Sprintf("%s.proto", table.Name.Snake().Lower()),
			Content:  []byte(prfile.Render(indent)),
		}

		pages = append(pages, expPage)
	}

	return pages, nil
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

	gaHTTP := modifiers.GoogleApiHttp{}
	gaFH := modifiers.GoogleAPIFieldBehavior{}

	prfile := presentation.AllocateFile(
		spec.Package,
		len(params.Schema.Enums),
		presentation.WithModifyProcedure(gaHTTP.ModifyProcedure),
		presentation.WithModifyField(gaFH.ModifyField),
	).SetOptions(options)

	for _, enum := range params.Schema.Enums {
		prfile.AddEnum(*enum.Name, enum.Values)
	}

	tablemsgMapper := e.newTableMapper(spec)

	for _, table := range params.Schema.Tables.List() {
		srv, err := e.buildService(tablemsgMapper, params.Schema.Driver, prfile, table, map[string]*exporter.ExportedPage{})
		if err != nil {
			return nil, fmt.Errorf("build service for table %q: %w", table.Name.Value, err)
		}
		if !srv.HasProcedures() {
			continue
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
	tablemsgMapper *tablemsg.Mapper,
	sourceDriver config.DatabaseDriver,
	prfile *presentation.File,
	table *schema.Table,
	enumPages map[string]*exporter.ExportedPage,
) (*presentation.Service, error) {
	procedureBuilders := map[presentation.ProcedureType]func(buildCtx *buildProcedureContext) error{
		presentation.ProcedureTypeList:   e.buildListProcedure,
		presentation.ProcedureTypeGet:    e.buildGetProcedure,
		presentation.ProcedureTypeDelete: e.buildDeleteProcedure,
		presentation.ProcedureTypeCreate: e.buildCreateProcedure,
		presentation.ProcedureTypePatch:  e.buildPatchProcedure,
	}

	tblmsg := tablemsgMapper.MapTable(prfile, table, func(col *schema.Column) string {
		return e.mapType(sourceDriver, col, prfile, enumPages)
	})

	srv := prfile.AddService(fmt.Sprintf("%sService", table.Name.Pascal()), tblmsg, 5)

	buildCtx := &buildProcedureContext{
		sourceDriver: sourceDriver,
		service:      srv,
		table:        table,
		tableMsg:     tblmsg,
		enumPages:    enumPages,
	}
	buildCtx.tableSingularName = buildCtx.service.TableMessage().Name()

	for procType, builder := range procedureBuilders {
		err := builder(buildCtx)
		if err != nil {
			return nil, fmt.Errorf("build procedure of %s: %w", string(procType), err)
		}
	}

	return srv, nil
}

func (e *Exporter) buildGetProcedure(
	buildCtx *buildProcedureContext,
) error {
	if buildCtx.table.PrimaryKey == nil {
		return nil
	}

	buildCtx.service.AddProcedureFn(
		"Get",
		presentation.ProcedureTypeGet,
		func(message *presentation.Message) {
			message.SetName(fmt.Sprintf("Get%sRequest", buildCtx.tableSingularName))

			for _, pkField := range buildCtx.tableMsg.PrimaryKey {
				message.CreateField(pkField.Name, func(field *presentation.Field) {
					if pkField.Repeated {
						field.AsRepeated()
					}
					field.SetType(pkField.Type).AsRequired()
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
	if buildCtx.table.PrimaryKey == nil {
		return nil
	}

	var err error

	buildCtx.service.AddProcedureFn(
		"Delete",
		presentation.ProcedureTypeDelete,
		func(message *presentation.Message) {
			message.SetName(fmt.Sprintf("Delete%sRequest", buildCtx.tableSingularName))

			for _, pkField := range buildCtx.tableMsg.PrimaryKey {
				message.CreateField(pkField.Name, func(field *presentation.Field) {
					if pkField.Repeated {
						field.AsRepeated()
					}
					field.SetType(pkField.Type)
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
	if buildCtx.table.PrimaryKey == nil {
		return nil
	}

	buildCtx.service.AddProcedureFn("Create", presentation.ProcedureTypeCreate,
		func(message *presentation.Message) {
			message.SetName(fmt.Sprintf("Create%sRequest", buildCtx.tableSingularName))

			for _, col := range buildCtx.table.Columns {
				if col.IsAutoincrement {
					continue
				}

				tableField, _ := message.Service().TableMessage().GetField(col.Name.Value)

				message.CreateField(tableField.Name, func(field *presentation.Field) {
					field.SetType(tableField.Type)

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

			for _, col := range buildCtx.table.Columns {
				tableField, _ := message.Service().TableMessage().GetField(col.Name.Value)

				message.CreateField(tableField.Name, func(field *presentation.Field) {
					field.SetType(tableField.Type)

					if !col.Nullable {
						field.AsRequired()
					}
				})
			}
		},
		func(message *presentation.Message) {
			message.
				SetName(fmt.Sprintf("Patch%sResponse", buildCtx.tableSingularName)).
				CreateField(buildCtx.table.Name.Pascal().Singular().Value, func(field *presentation.Field) {
					field.SetType(buildCtx.service.TableMessage().Name())
				})
		},
	)

	return nil
}

func (e *Exporter) mapType(
	sourceDriver config.DatabaseDriver,
	column *schema.Column,
	file *presentation.File,
	enumPages map[string]*exporter.ExportedPage,
) string {
	if column.Enum != nil {
		enumPage, enumPageExists := enumPages[column.Enum.Name.Value]
		if enumPageExists {
			file.AddImport(enumPage.FileName)
		}

		return column.Enum.Name.Pascal().Value
	}

	goType := sqltype.MapGoType(sourceDriver, column.Type)

	switch goType {
	case golang.TypeInt, golang.TypeInt16, golang.TypeInt64:
		return "int64"
	case golang.TypeFloat64:
		return "double"
	case golang.TypeFloat32:
		return "double"
	case golang.TypeBool:
		return "bool"
	case golang.TypeTimeTime:
		file.AddImport("google/protobuf/timestamp.proto")

		return "google.protobuf.Timestamp"
	default:
		return "string"
	}
}
