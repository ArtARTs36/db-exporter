package grpccrud

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/infrastructure/sqltype"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/golang"
	"github.com/artarts36/db-exporter/internal/shared/proto"
	"github.com/artarts36/gds"
)

type Exporter struct{}

type buildProcedureContext struct {
	sourceDriver config.DatabaseDriver

	prfile            *proto.File
	table             *schema.Table
	tableMsg          *proto.Message
	tableSingularName string
	enumPages         map[string]*exporter.ExportedPage
}

func NewCrudExporter() *Exporter {
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
	indent := proto.NewIndent(spec.Indent)

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

	for _, table := range params.Schema.Tables.List() {
		prfile := &proto.File{
			Package:  spec.Package,
			Services: make([]*proto.Service, 0, 1),
			Messages: make([]*proto.Message, 0, params.Schema.Tables.Len()),
			Imports:  gds.NewSet[string](),
			Options:  options,
		}

		srv, messages := e.buildService(params.Schema.Driver, prfile, table, enumPages)

		if len(srv.Procedures) == 0 {
			continue
		}

		prfile.Services = append(prfile.Services, srv)
		prfile.Messages = append(prfile.Messages, messages...)

		expPage := &exporter.ExportedPage{
			FileName: fmt.Sprintf("%s.proto", table.Name.Lower().Lower()),
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

	for _, table := range params.Schema.Tables.List() {
		srv, messages := e.buildService(params.Schema.Driver, prfile, table, map[string]*exporter.ExportedPage{})

		if len(srv.Procedures) == 0 {
			continue
		}

		prfile.Services = append(prfile.Services, srv)
		prfile.Messages = append(prfile.Messages, messages...)
	}

	expPage := &exporter.ExportedPage{
		FileName: "services.proto",
		Content:  []byte(prfile.Render(proto.NewIndent(spec.Indent))),
	}

	return []*exporter.ExportedPage{
		expPage,
	}, nil
}

func (e *Exporter) buildService(
	sourceDriver config.DatabaseDriver,
	prfile *proto.File,
	table *schema.Table,
	enumPages map[string]*exporter.ExportedPage,
) (*proto.Service, []*proto.Message) {
	procedureBuilders := []func(buildCtx *buildProcedureContext) (
		*proto.ServiceProcedure,
		[]*proto.Message,
	){
		e.buildListProcedure,
		e.buildGetProcedure,
		e.buildDeleteProcedure,
		e.buildCreateProcedure,
		e.buildPatchProcedure,
	}

	messages := make([]*proto.Message, 0, 1)

	srv := &proto.Service{
		Name: fmt.Sprintf("%sService", table.Name.Pascal()),
	}

	buildCtx := &buildProcedureContext{
		sourceDriver: sourceDriver,
		prfile:       prfile,
		table:        table,
		tableMsg: &proto.Message{
			Name:   table.Name.Pascal().Singular().Value,
			Fields: make([]*proto.Field, 0, len(table.Columns)),
		},
		enumPages: enumPages,
	}
	buildCtx.tableSingularName = buildCtx.tableMsg.Name

	messages = append(messages, buildCtx.tableMsg)

	id := 1

	for _, column := range table.Columns {
		buildCtx.tableMsg.Fields = append(buildCtx.tableMsg.Fields, &proto.Field{
			Name: column.Name.Lower().Value,
			Type: e.mapType(buildCtx.sourceDriver, column, prfile.Imports, buildCtx.enumPages),
			ID:   id,
		})

		id++
	}

	for _, builder := range procedureBuilders {
		procedure, msgs := builder(buildCtx)
		if procedure == nil {
			continue
		}

		srv.Procedures = append(srv.Procedures, procedure)
		messages = append(messages, msgs...)
	}

	return srv, messages
}

func (e *Exporter) buildGetProcedure(
	buildCtx *buildProcedureContext,
) (*proto.ServiceProcedure, []*proto.Message) {
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

		getReqMsg.Fields = append(getReqMsg.Fields, &proto.Field{
			Name: col.Name.Lower().Value,
			Type: e.mapType(buildCtx.sourceDriver, col, buildCtx.prfile.Imports, buildCtx.enumPages),
			ID:   id,
		})

		id++
	}

	getRespMsg := &proto.Message{
		Name: fmt.Sprintf("Get%sResponse", buildCtx.tableSingularName),
		Fields: []*proto.Field{
			{
				Name: buildCtx.table.Name.Pascal().Singular().Value,
				Type: buildCtx.tableMsg.Name,
				ID:   1,
			},
		},
	}

	return &proto.ServiceProcedure{
		Name:    "Get",
		Param:   getReqMsg.Name,
		Returns: getRespMsg.Name,
	}, []*proto.Message{getReqMsg, getRespMsg}
}

func (e *Exporter) buildListProcedure(
	buildCtx *buildProcedureContext,
) (*proto.ServiceProcedure, []*proto.Message) {
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
				Type:     buildCtx.tableMsg.Name,
				Name:     "items",
				ID:       1,
			},
		},
	}

	return &proto.ServiceProcedure{
		Name:    "List",
		Param:   getReqMsg.Name,
		Returns: respMsg.Name,
	}, []*proto.Message{getReqMsg, respMsg}
}

func (e *Exporter) buildDeleteProcedure(
	buildCtx *buildProcedureContext,
) (*proto.ServiceProcedure, []*proto.Message) {
	if buildCtx.table.PrimaryKey == nil {
		return nil, nil
	}

	deleteReqMsg := &proto.Message{
		Name:   fmt.Sprintf("Delete%sRequest", buildCtx.tableSingularName),
		Fields: make([]*proto.Field, 0),
	}

	deleteRespMsg := &proto.Message{
		Name: fmt.Sprintf("Delete%sResponse", buildCtx.tableSingularName),
	}

	id := 1

	for _, col := range buildCtx.table.Columns {
		if !col.IsPrimaryKey() {
			continue
		}

		deleteReqMsg.Fields = append(deleteReqMsg.Fields, &proto.Field{
			Name: col.Name.Lower().Value,
			Type: e.mapType(buildCtx.sourceDriver, col, buildCtx.prfile.Imports, buildCtx.enumPages),
			ID:   id,
		})

		id++
	}

	return &proto.ServiceProcedure{
		Name:    "Delete",
		Param:   deleteReqMsg.Name,
		Returns: deleteRespMsg.Name,
	}, []*proto.Message{deleteReqMsg, deleteRespMsg}
}

func (e *Exporter) buildCreateProcedure(
	buildCtx *buildProcedureContext,
) (*proto.ServiceProcedure, []*proto.Message) {
	if buildCtx.table.PrimaryKey == nil {
		return nil, nil
	}

	createReqMsg := &proto.Message{
		Name:   fmt.Sprintf("Create%sRequest", buildCtx.tableSingularName),
		Fields: make([]*proto.Field, 0, len(buildCtx.table.Columns)),
	}

	createRespMsg := &proto.Message{
		Name: fmt.Sprintf("Create%sResponse", buildCtx.tableSingularName),
		Fields: []*proto.Field{
			{
				Name: buildCtx.table.Name.Pascal().Singular().Value,
				Type: buildCtx.tableMsg.Name,
				ID:   1,
			},
		},
	}

	id := 1

	for _, col := range buildCtx.table.Columns {
		if col.IsAutoincrement {
			continue
		}

		createReqMsg.Fields = append(createReqMsg.Fields, &proto.Field{
			Name: col.Name.Lower().Value,
			Type: e.mapType(buildCtx.sourceDriver, col, buildCtx.prfile.Imports, buildCtx.enumPages),
			ID:   id,
		})

		id++
	}

	return &proto.ServiceProcedure{
		Name:    "Create",
		Param:   createReqMsg.Name,
		Returns: createRespMsg.Name,
	}, []*proto.Message{createReqMsg, createRespMsg}
}

func (e *Exporter) buildPatchProcedure(
	buildCtx *buildProcedureContext,
) (*proto.ServiceProcedure, []*proto.Message) {
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
				Type: buildCtx.tableMsg.Name,
				ID:   1,
			},
		},
	}

	id := 1

	for _, col := range buildCtx.table.Columns {
		if col.IsAutoincrement {
			continue
		}

		patchReqMsg.Fields = append(patchReqMsg.Fields, &proto.Field{
			Name: col.Name.Lower().Value,
			Type: e.mapType(buildCtx.sourceDriver, col, buildCtx.prfile.Imports, buildCtx.enumPages),
			ID:   id,
		})

		id++
	}

	return &proto.ServiceProcedure{
		Name:    "Patch",
		Param:   patchReqMsg.Name,
		Returns: patchRespMsg.Name,
	}, []*proto.Message{patchReqMsg, patchRespMsg}
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
