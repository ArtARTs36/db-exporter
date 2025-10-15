package grpccrud

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/fieldmap"
	service "github.com/artarts36/db-exporter/internal/exporter/grpc-crud/service"
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/tablemsg"
	"github.com/artarts36/db-exporter/internal/infrastructure/sqltype"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/golang"
	"github.com/artarts36/db-exporter/internal/shared/indentx"
	"github.com/artarts36/db-exporter/internal/shared/proto"
	"github.com/artarts36/gds"
)

type Exporter struct{}

type buildProcedureContext struct {
	sourceDriver config.DatabaseDriver

	prfile            *proto.File
	table             *schema.Table
	tableMsg          *tablemsg.Message
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
	procModifier := service.SelectProcedureModifier(spec)

	enumPages := map[string]*exporter.ExportedPage{}
	indent := indentx.NewIndent(spec.Indent)

	for _, enum := range params.Schema.Enums {
		prfile := &proto.File{
			Package: spec.Package,
			Options: options,
			Enums:   []*proto.Enum{proto.NewEnumWithValues(enum.Name, enum.Values)},
		}

		expPage := &exporter.ExportedPage{
			FileName: fmt.Sprintf("%s_enum.proto", enum.Name.Value),
			Content:  []byte(prfile.Render(indent)),
		}

		enumPages[enum.Name.Value] = expPage
		pages = append(pages, expPage)
	}

	tablemsgMapper := e.newTableMapper(spec)

	for _, table := range params.Schema.Tables.List() {
		prfile := &proto.File{
			Package:  spec.Package,
			Services: make([]*proto.Service, 0, 1),
			Messages: make([]*proto.Message, 0, params.Schema.Tables.Len()),
			Imports:  gds.NewSet[string](),
			Options:  options,
		}

		srv, err := e.buildService(tablemsgMapper, params.Schema.Driver, prfile, table, enumPages, procModifier)
		if err != nil {
			return nil, fmt.Errorf("build service for table %q: %w", table.Name, err)
		}
		if len(srv.Procedures) == 0 {
			continue
		}

		prfile.Services = append(prfile.Services, srv.ToProto())
		prfile.Messages = append(prfile.Messages, srv.Messages...)

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
	procModifier := service.SelectProcedureModifier(spec)

	prfile := &proto.File{
		Package:  spec.Package,
		Services: make([]*proto.Service, 0, params.Schema.Tables.Len()),
		Messages: make([]*proto.Message, 0, params.Schema.Tables.Len()),
		Imports:  gds.NewSet[string](),
		Options:  options,
		Enums:    make([]*proto.Enum, 0, len(params.Schema.Enums)),
	}

	for _, enum := range params.Schema.Enums {
		prfile.Enums = append(prfile.Enums, proto.NewEnumWithValues(enum.Name, enum.Values))
	}

	tablemsgMapper := e.newTableMapper(spec)

	for _, table := range params.Schema.Tables.List() {
		srv, err := e.buildService(tablemsgMapper, params.Schema.Driver, prfile, table, map[string]*exporter.ExportedPage{}, procModifier)
		if err != nil {
			return nil, fmt.Errorf("build service for table %q: %w", table.Name.Value, err)
		}
		if len(srv.Procedures) == 0 {
			continue
		}

		prfile.Services = append(prfile.Services, srv.ToProto())
		prfile.Messages = append(prfile.Messages, srv.Messages...)
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
	prfile *proto.File,
	table *schema.Table,
	enumPages map[string]*exporter.ExportedPage,
	createProcModifier service.ProcedureModifierFactory,
) (*service.Service, error) {
	procedureBuilders := map[service.ProcedureType]func(buildCtx *buildProcedureContext) (*service.Procedure, error){
		service.ProcedureTypeList:   e.buildListProcedure,
		service.ProcedureTypeGet:    e.buildGetProcedure,
		service.ProcedureTypeDelete: e.buildDeleteProcedure,
		service.ProcedureTypeCreate: e.buildCreateProcedure,
		service.ProcedureTypePatch:  e.buildPatchProcedure,
	}

	srv := &service.Service{
		Name: fmt.Sprintf("%sService", table.Name.Pascal()),
	}

	buildCtx := &buildProcedureContext{
		sourceDriver: sourceDriver,
		prfile:       prfile,
		table:        table,
		tableMsg: tablemsgMapper.MapTable(prfile, table, func(col *schema.Column) string {
			return e.mapType(sourceDriver, col, prfile.Imports, enumPages)
		}),
		enumPages: enumPages,
	}
	buildCtx.tableSingularName = buildCtx.tableMsg.Proto.Name

	srv.Messages = append(srv.Messages, buildCtx.tableMsg.Proto)

	procModifier := createProcModifier(prfile, srv, buildCtx.tableMsg)

	for procType, builder := range procedureBuilders {
		proc, err := builder(buildCtx)
		if err != nil {
			return nil, fmt.Errorf("build procedure of %s: %w", string(procType), err)
		}
		if proc == nil {
			continue
		}

		procModifier(proc)

		srv.Procedures = append(srv.Procedures, proc)
		srv.Messages = append(srv.Messages, proc.Request, proc.Response)
	}

	return srv, nil
}

func (e *Exporter) buildGetProcedure(
	buildCtx *buildProcedureContext,
) (*service.Procedure, error) {
	if buildCtx.table.PrimaryKey == nil {
		return nil, nil
	}

	getReqMsg := &proto.Message{
		Name:   fmt.Sprintf("Get%sRequest", buildCtx.tableSingularName),
		Fields: make([]*proto.Field, 0),
	}

	id := 1

	for _, col := range buildCtx.table.Columns {
		if !col.IsPrimaryKey() {
			continue
		}

		field, err := buildCtx.tableMsg.CloneField(col.Name.Value)
		if err != nil {
			return nil, err
		}

		getReqMsg.Fields = append(getReqMsg.Fields, field)

		id++
	}

	return &service.Procedure{
		Name:    "Get",
		Type:    service.ProcedureTypeGet,
		Request: getReqMsg,
		Response: &proto.Message{
			Name: fmt.Sprintf("Get%sResponse", buildCtx.tableSingularName),
			Fields: []*proto.Field{
				{
					Name: buildCtx.table.Name.Pascal().Singular().Value,
					Type: buildCtx.tableMsg.Proto.Name,
					ID:   1,
				},
			},
		},
	}, nil
}

func (e *Exporter) buildListProcedure(
	buildCtx *buildProcedureContext,
) (*service.Procedure, error) {
	if buildCtx.table.PrimaryKey == nil {
		return nil, nil
	}

	getReqMsg := &proto.Message{
		Name:   fmt.Sprintf("List%sRequest", buildCtx.tableSingularName),
		Fields: make([]*proto.Field, 0),
	}

	respMsg := &proto.Message{
		Name: fmt.Sprintf("List%sResponse", buildCtx.tableSingularName),
		Fields: []*proto.Field{
			{
				Repeated: true,
				Type:     buildCtx.tableMsg.Proto.Name,
				Name:     "items",
				ID:       1,
			},
		},
	}

	return &service.Procedure{
		Name:     "List",
		Type:     service.ProcedureTypeList,
		Request:  getReqMsg,
		Response: respMsg,
	}, nil
}

func (e *Exporter) buildDeleteProcedure(
	buildCtx *buildProcedureContext,
) (*service.Procedure, error) {
	if buildCtx.table.PrimaryKey == nil {
		return nil, nil
	}

	deleteReqMsg := &proto.Message{
		Name:   fmt.Sprintf("Delete%sRequest", buildCtx.tableSingularName),
		Fields: make([]*proto.Field, 0),
	}

	id := 1

	for _, col := range buildCtx.table.Columns {
		if !col.IsPrimaryKey() {
			continue
		}

		field, err := buildCtx.tableMsg.CloneField(col.Name.Value)
		if err != nil {
			return nil, err
		}

		deleteReqMsg.Fields = append(deleteReqMsg.Fields, field)

		id++
	}

	return &service.Procedure{
		Name:    "Delete",
		Type:    service.ProcedureTypeDelete,
		Request: deleteReqMsg,
		Response: &proto.Message{
			Name: fmt.Sprintf("Delete%sResponse", buildCtx.tableSingularName),
		},
	}, nil
}

func (e *Exporter) buildCreateProcedure(
	buildCtx *buildProcedureContext,
) (*service.Procedure, error) {
	if buildCtx.table.PrimaryKey == nil {
		return nil, nil
	}

	createReqMsg := &proto.Message{
		Name:   fmt.Sprintf("Create%sRequest", buildCtx.tableSingularName),
		Fields: make([]*proto.Field, 0, len(buildCtx.table.Columns)),
	}

	id := 1

	for _, col := range buildCtx.table.Columns {
		if col.IsAutoincrement {
			continue
		}

		field, err := buildCtx.tableMsg.CloneField(col.Name.Value)
		if err != nil {
			return nil, err
		}

		createReqMsg.Fields = append(createReqMsg.Fields, field)

		id++
	}

	return &service.Procedure{
		Name:    "Create",
		Type:    service.ProcedureTypeCreate,
		Request: createReqMsg,
		Response: &proto.Message{
			Name: fmt.Sprintf("Create%sResponse", buildCtx.tableSingularName),
			Fields: []*proto.Field{
				{
					Name: buildCtx.table.Name.Pascal().Singular().Value,
					Type: buildCtx.tableMsg.Proto.Name,
					ID:   1,
				},
			},
		},
	}, nil
}

func (e *Exporter) buildPatchProcedure(
	buildCtx *buildProcedureContext,
) (*service.Procedure, error) {
	if buildCtx.table.PrimaryKey == nil {
		return nil, nil
	}

	patchReqMsg := &proto.Message{
		Name:   fmt.Sprintf("Patch%sRequest", buildCtx.tableSingularName),
		Fields: make([]*proto.Field, 0, len(buildCtx.tableSingularName)),
	}

	patchRespMsg := &proto.Message{
		Name: fmt.Sprintf("Patch%sResponse", buildCtx.tableSingularName),
		Fields: []*proto.Field{
			{
				Name: buildCtx.table.Name.Pascal().Singular().Value,
				Type: buildCtx.tableMsg.Proto.Name,
				ID:   1,
			},
		},
	}

	id := 1

	for _, col := range buildCtx.table.Columns {
		if col.IsAutoincrement {
			continue
		}

		field, err := buildCtx.tableMsg.CloneField(col.Name.Value)
		if err != nil {
			return nil, err
		}

		patchReqMsg.Fields = append(patchReqMsg.Fields, field)

		id++
	}

	return &service.Procedure{
		Name:     "Patch",
		Type:     service.ProcedureTypePatch,
		Request:  patchReqMsg,
		Response: patchRespMsg,
	}, nil
}

func (e *Exporter) mapType(
	sourceDriver config.DatabaseDriver,
	column *schema.Column,
	imports *gds.Set[string],
	enumPages map[string]*exporter.ExportedPage,
) string {
	if column.Enum != nil {
		enumPage, enumPageExists := enumPages[column.Enum.Name.Value]
		if enumPageExists {
			imports.Add(enumPage.FileName)
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
		imports.Add("google/protobuf/timestamp.proto")

		return "google.protobuf.Timestamp"
	default:
		return "string"
	}
}
