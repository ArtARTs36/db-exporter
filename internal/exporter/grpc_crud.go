package exporter

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/ds"
	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/proto"
	"github.com/artarts36/db-exporter/internal/template"
)

const GrpcCrudExporterName = "grpc-crud"

type GrpcCrudExporter struct {
	renderer *template.Renderer
}

func NewGrpcCrudExporter(renderer *template.Renderer) *GrpcCrudExporter {
	return &GrpcCrudExporter{
		renderer: renderer,
	}
}

func (e *GrpcCrudExporter) ExportPerFile(
	_ context.Context,
	sc *schema.Schema,
	params *ExportParams,
) ([]*ExportedPage, error) {
	pages := make([]*ExportedPage, 0, sc.Tables.Len())

	for _, table := range sc.Tables.List() {
		prfile := &proto.File{
			Package:  params.Package,
			Services: make([]*proto.Service, 0, sc.Tables.Len()),
			Messages: make([]*proto.Message, 0, sc.Tables.Len()),
			Imports:  ds.NewSet(),
		}

		srv, messages := e.buildService(prfile, table)

		if len(srv.Procedures) == 0 {
			continue
		}

		prfile.Services = append(prfile.Services, srv)
		prfile.Messages = append(prfile.Messages, messages...)

		expPage, err := render(
			e.renderer,
			"grpc-crud/grpc.proto",
			fmt.Sprintf("%s.proto", table.Name.Lower().Lower()),
			map[string]stick.Value{
				"file": prfile,
			},
		)
		if err != nil {
			return nil, err
		}

		pages = append(pages, expPage)
	}

	return pages, nil
}

func (e *GrpcCrudExporter) Export(
	_ context.Context,
	sc *schema.Schema,
	params *ExportParams,
) ([]*ExportedPage, error) {
	prfile := &proto.File{
		Package:  params.Package,
		Services: make([]*proto.Service, 0, sc.Tables.Len()),
		Messages: make([]*proto.Message, 0, sc.Tables.Len()),
		Imports:  ds.NewSet(),
	}

	for _, table := range sc.Tables.List() {
		srv, messages := e.buildService(prfile, table)

		if len(srv.Procedures) == 0 {
			continue
		}

		prfile.Services = append(prfile.Services, srv)
		prfile.Messages = append(prfile.Messages, messages...)
	}

	expPage, err := render(e.renderer, "grpc-crud/grpc.proto", "services.proto", map[string]stick.Value{
		"file": prfile,
	})
	if err != nil {
		return nil, err
	}

	return []*ExportedPage{
		expPage,
	}, nil
}

func (e *GrpcCrudExporter) buildService(prfile *proto.File, table *schema.Table) (*proto.Service, []*proto.Message) {
	procedureBuilders := []func(prfile *proto.File, table *schema.Table, tableMsg *proto.Message) (
		*proto.ServiceProcedure,
		[]*proto.Message,
	){
		e.buildGetListProcedure,
		e.buildGetProcedure,
		e.buildDeleteProcedure,
		e.buildCreateProcedure,
		e.buildPatchProcedure,
	}

	messages := make([]*proto.Message, 0, 1)

	srv := &proto.Service{
		Name: fmt.Sprintf("%sService", table.Name.Pascal()),
	}

	tableMsg := &proto.Message{
		Name:   table.Name.Pascal().Singular().Value,
		Fields: make([]*proto.Field, 0, len(table.Columns)),
	}

	messages = append(messages, tableMsg)

	id := 1

	for _, column := range table.Columns {
		tableMsg.Fields = append(tableMsg.Fields, &proto.Field{
			Name: column.Name.Lower().Value,
			Type: e.mapType(column, prfile.Imports),
			ID:   id,
		})

		id++
	}

	for _, builder := range procedureBuilders {
		procedure, msgs := builder(prfile, table, tableMsg)
		if procedure == nil {
			continue
		}

		srv.Procedures = append(srv.Procedures, procedure)
		messages = append(messages, msgs...)
	}

	return srv, messages
}

func (e *GrpcCrudExporter) buildGetProcedure(
	prfile *proto.File,
	table *schema.Table,
	tableMsg *proto.Message,
) (*proto.ServiceProcedure, []*proto.Message) {
	if table.PrimaryKey == nil {
		return nil, nil
	}

	getReqMsg := &proto.Message{
		Name:   fmt.Sprintf("Get%sRequest", table.Name.Pascal().Singular().Value),
		Fields: make([]*proto.Field, 0),
	}

	id := 1

	for _, col := range table.Columns {
		if !col.IsPrimaryKey() {
			continue
		}

		getReqMsg.Fields = append(getReqMsg.Fields, &proto.Field{
			Name: col.Name.Lower().Value,
			Type: e.mapType(col, prfile.Imports),
			ID:   id,
		})

		id++
	}

	return &proto.ServiceProcedure{
		Name:    fmt.Sprintf("Get%s", table.Name.Pascal().Singular()),
		Param:   getReqMsg.Name,
		Returns: tableMsg.Name,
	}, []*proto.Message{getReqMsg}
}

func (e *GrpcCrudExporter) buildGetListProcedure(
	_ *proto.File,
	table *schema.Table,
	tableMsg *proto.Message,
) (*proto.ServiceProcedure, []*proto.Message) {
	if table.PrimaryKey == nil {
		return nil, nil
	}

	getReqMsg := &proto.Message{
		Name:   fmt.Sprintf("Get%sListRequest", table.Name.Pascal().Singular().Value),
		Fields: make([]*proto.Field, 0),
	}

	respMsg := &proto.Message{
		Name: fmt.Sprintf("Get%sResponse", table.Name.Pascal().Singular().Value),
		Fields: []*proto.Field{
			{
				Repeated: true,
				Type:     tableMsg.Name,
				Name:     "items",
				ID:       1,
			},
		},
	}

	return &proto.ServiceProcedure{
		Name:    fmt.Sprintf("Get%sList", table.Name.Pascal().Singular()),
		Param:   getReqMsg.Name,
		Returns: tableMsg.Name,
	}, []*proto.Message{getReqMsg, respMsg}
}

func (e *GrpcCrudExporter) buildDeleteProcedure(
	prfile *proto.File,
	table *schema.Table,
	_ *proto.Message,
) (*proto.ServiceProcedure, []*proto.Message) {
	if table.PrimaryKey == nil {
		return nil, nil
	}

	deleteReqMsg := &proto.Message{
		Name:   fmt.Sprintf("Delete%sRequest", table.Name.Pascal().Singular().Value),
		Fields: make([]*proto.Field, 0),
	}

	deleteRespMsg := &proto.Message{
		Name: fmt.Sprintf("Delete%sResponse", table.Name.Pascal().Singular()),
	}

	id := 1

	for _, col := range table.Columns {
		if !col.IsPrimaryKey() {
			continue
		}

		deleteReqMsg.Fields = append(deleteReqMsg.Fields, &proto.Field{
			Name: col.Name.Lower().Value,
			Type: e.mapType(col, prfile.Imports),
			ID:   id,
		})

		id++
	}

	return &proto.ServiceProcedure{
		Name:    fmt.Sprintf("Delete%s", table.Name.Pascal().Singular()),
		Param:   deleteReqMsg.Name,
		Returns: deleteRespMsg.Name,
	}, []*proto.Message{deleteReqMsg, deleteRespMsg}
}

func (e *GrpcCrudExporter) buildCreateProcedure(
	prfile *proto.File,
	table *schema.Table,
	tableMsg *proto.Message,
) (*proto.ServiceProcedure, []*proto.Message) {
	if table.PrimaryKey == nil {
		return nil, nil
	}

	createReqMsg := &proto.Message{
		Name:   fmt.Sprintf("Create%sRequest", table.Name.Pascal().Singular().Value),
		Fields: make([]*proto.Field, 0, len(table.Columns)),
	}

	id := 1

	for _, col := range table.Columns {
		if col.IsPrimaryKey() {
			continue
		}

		createReqMsg.Fields = append(createReqMsg.Fields, &proto.Field{
			Name: col.Name.Lower().Value,
			Type: e.mapType(col, prfile.Imports),
			ID:   id,
		})

		id++
	}

	return &proto.ServiceProcedure{
		Name:    fmt.Sprintf("Create%s", table.Name.Pascal().Singular()),
		Param:   createReqMsg.Name,
		Returns: tableMsg.Name,
	}, []*proto.Message{createReqMsg}
}

func (e *GrpcCrudExporter) buildPatchProcedure(
	prfile *proto.File,
	table *schema.Table,
	tableMsg *proto.Message,
) (*proto.ServiceProcedure, []*proto.Message) {
	if table.PrimaryKey == nil {
		return nil, nil
	}

	patchReqMsg := &proto.Message{
		Name:   fmt.Sprintf("Patch%sRequest", table.Name.Pascal().Singular().Value),
		Fields: make([]*proto.Field, 0, len(table.Columns)),
	}

	id := 1

	for _, col := range table.Columns {
		patchReqMsg.Fields = append(patchReqMsg.Fields, &proto.Field{
			Name: col.Name.Lower().Value,
			Type: e.mapType(col, prfile.Imports),
			ID:   id,
		})

		id++
	}

	return &proto.ServiceProcedure{
		Name:    fmt.Sprintf("Patch%s", table.Name.Pascal().Singular()),
		Param:   patchReqMsg.Name,
		Returns: tableMsg.Name,
	}, []*proto.Message{patchReqMsg}
}

func (e *GrpcCrudExporter) mapType(column *schema.Column, imports *ds.Set) string {
	switch column.PreparedType { //nolint: exhaustive // not need
	case schema.ColumnTypeInteger:
		return "int64"
	case schema.ColumnTypeFloat64:
		return "double"
	case schema.ColumnTypeFloat32:
		return "double"
	case schema.ColumnTypeBoolean:
		return "bool"
	case schema.ColumnTypeTimestamp:
		imports.Add("google/protobuf/timestamp.proto")

		return "google.protobuf.Timestamp"
	default:
		return "string"
	}
}
