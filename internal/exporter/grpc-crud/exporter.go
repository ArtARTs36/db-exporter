package grpccrud

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/modifiers"
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/paginator"
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/iox"
	"github.com/artarts36/db-exporter/internal/shared/proto"
	"github.com/artarts36/gds"
)

// Exporter based on https://google.aip.dev/121
type Exporter struct{}

type buildProcedureContext struct {
	spec *Specification

	sourceDriver schema.DatabaseDriver

	service *presentation.Service

	paginator paginator.Paginator

	tableSingularName string
	tablePluralName   string
}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) createPaginator(spec *Specification) paginator.Paginator {
	if spec.Pagination == paginationTypeOffset {
		return &paginator.Offset{}
	}

	if spec.Pagination == paginationTypeNone {
		return &paginator.None{}
	}

	return &paginator.Token{}
}

func (e *Exporter) ExportPerFile(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*Specification)
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

		prfile := pkg.CreateFile(fmt.Sprintf("%s_enum.proto", enum.Name.Value)).SetOptions(options)
		prfile.AddEnum(*enum.Name, enum.Values)

		expPage := &exporter.ExportedPage{
			FileName: prfile.Name(),
			Content:  []byte(prfile.Render(indent)),
		}

		pages = append(pages, expPage)
	}

	pager := e.createPaginator(spec)

	for _, table := range params.Schema.Tables.List() {
		prfile := pkg.CreateFile(fmt.Sprintf("%s.proto", table.Name.Snake().Lower())).SetOptions(options)

		err := e.buildService(spec, params.Schema.Driver, prfile, table, pager)
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

func (e *Exporter) newPackage(spec *Specification) *presentation.Package {
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
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*Specification)
	if !ok {
		return nil, fmt.Errorf("invalid spec, expected Specification, got %T", params.Spec)
	}

	options := proto.PrepareOptions(spec.Options)

	pkg := e.newPackage(spec)
	pager := e.createPaginator(spec)

	prfile := pkg.CreateFile("services.proto").SetOptions(options)

	for _, enum := range params.Schema.Enums {
		if enum.Used == 0 {
			continue
		}

		prfile.AddEnum(*enum.Name, enum.Values)
	}

	for _, table := range params.Schema.Tables.List() {
		err := e.buildService(spec, params.Schema.Driver, prfile, table, pager)
		if err != nil {
			return nil, fmt.Errorf("build service for table %q: %w", table.Name.Value, err)
		}
	}

	expPage := &exporter.ExportedPage{
		FileName: "services.proto",
		Content:  []byte(prfile.Render(iox.NewIndent(spec.Indent))),
	}

	return []*exporter.ExportedPage{
		expPage,
	}, nil
}

func (e *Exporter) buildService(
	spec *Specification,
	sourceDriver schema.DatabaseDriver,
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
			Type:  presentation.ProcedureTypeUndelete,
			Build: e.buildUndeleteProcedure,
		},
		{
			Type:  presentation.ProcedureTypeCreate,
			Build: e.buildCreateProcedure,
		},
		{
			Type:  presentation.ProcedureTypeUpdate,
			Build: e.buildUpdateProcedure,
		},
	}

	mapColumnType := func(col *schema.Column) string {
		return e.mapType(sourceDriver, col, prfile)
	}

	srv := prfile.AddService(
		table,
		func(message *presentation.TableMessage) {
			e.mapTableMessage(message, table, mapColumnType)
		},
	)

	buildCtx := &buildProcedureContext{
		spec:         spec,
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

func (e *Exporter) mapTableMessage(
	message *presentation.TableMessage,
	table *schema.Table,
	mapColumnType func(col *schema.Column) string,
) {
	for _, column := range table.Columns {
		creator := func(field *presentation.Field) {
			field.SetType(mapColumnType(column))

			if !column.Nullable {
				field.AsRequired()
			}

			field.SetColumn(column)

			if column.Comment.IsNotEmpty() {
				field.SetTopComment(column.Comment.WithSuffix(".").Value)
			}

			if e.columnAutofilled(column) {
				field.AsAutofilled()
			}
		}

		fieldName := column.Name.Snake().Lower().Value

		if column.IsPrimaryKey() {
			message.CreatePrimaryKeyField(fieldName, column.Name.Value, creator)
		} else {
			message.CreateField(fieldName, column.Name.Value, creator)
		}
	}
}

func (e *Exporter) buildGetProcedure(
	buildCtx *buildProcedureContext,
) error {
	if buildCtx.service.TableMessage().PrimaryKey == nil {
		return nil
	}

	buildCtx.service.AddProcedure(
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
	)

	return nil
}

// https://google.aip.dev/164
func (e *Exporter) buildUndeleteProcedure(
	buildCtx *buildProcedureContext,
) error {
	if buildCtx.service.TableMessage().PrimaryKey == nil {
		return nil
	}

	if !buildCtx.service.TableMessage().Table().SupportsSoftDelete() {
		return nil
	}

	buildCtx.service.AddProcedure(
		"Undelete",
		presentation.ProcedureTypeUndelete,
		func(message *presentation.Message) {
			message.SetName(fmt.Sprintf("Undelete%sRequest", buildCtx.tableSingularName))

			for _, pkField := range buildCtx.service.TableMessage().PrimaryKey {
				message.CreateField(pkField.Name(), func(field *presentation.Field) {
					field.CopyType(pkField).AsRequired()
				})
			}
		},
	)

	return nil
}

// https://google.aip.dev/132
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

			if message.Service().TableMessage().Table().SupportsSoftDelete() {
				message.CreateField("show_deleted", func(field *presentation.Field) {
					field.SetType("bool")
					field.SetTopComment("If set to `true`, soft-deleted resources will be returned alongside active resources.")
				})
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

// https://google.aip.dev/135
func (e *Exporter) buildDeleteProcedure(
	buildCtx *buildProcedureContext,
) error {
	if buildCtx.service.TableMessage().Table().PrimaryKey == nil {
		return nil
	}

	var err error

	requestFn := func(message *presentation.Message) {
		message.SetName(fmt.Sprintf("Delete%sRequest", buildCtx.tableSingularName))

		for _, pkField := range buildCtx.service.TableMessage().PrimaryKey {
			message.CreateField(pkField.Name(), func(field *presentation.Field) {
				field.CopyType(pkField)
				field.AsRequired()
			})
		}
	}

	if buildCtx.spec.RPC.Delete.Returns == deleteReturnsWrapper {
		buildCtx.service.AddProcedureFn(
			"Delete",
			presentation.ProcedureTypeDelete,
			requestFn,
			func(message *presentation.Message) {
				message.SetName(fmt.Sprintf("Delete%sResponse", buildCtx.tableSingularName))
			},
		)
	} else {
		buildCtx.service.File().AddImport("google/protobuf/empty.proto")

		buildCtx.service.AddProcedureWithResponseName(
			"Delete",
			presentation.ProcedureTypeDelete,
			requestFn,
			"google.protobuf.Empty",
		)
	}

	return err
}

// https://google.aip.dev/133
func (e *Exporter) buildCreateProcedure(
	buildCtx *buildProcedureContext,
) error {
	buildCtx.service.AddProcedure("Create", presentation.ProcedureTypeCreate,
		func(message *presentation.Message) {
			message.SetName(fmt.Sprintf("Create%sRequest", buildCtx.tableSingularName))

			message.CreateField(buildCtx.service.TableMessage().SingularNameForField(), func(fld *presentation.Field) {
				fld.
					SetType(buildCtx.service.TableMessage().Name()).
					AsRequired()
			})
		},
	)

	return nil
}

// https://google.aip.dev/134
func (e *Exporter) buildUpdateProcedure(
	buildCtx *buildProcedureContext,
) error {
	if buildCtx.service.TableMessage().Table().PrimaryKey == nil {
		return nil
	}

	buildCtx.service.File().AddImport("google/protobuf/field_mask.proto")

	buildCtx.service.AddProcedure("Update", presentation.ProcedureTypeUpdate,
		func(message *presentation.Message) {
			message.SetName(fmt.Sprintf("Update%sRequest", buildCtx.tableSingularName))

			message.CreateField(buildCtx.service.TableMessage().SingularNameForField(), func(fld *presentation.Field) {
				fld.
					SetType(buildCtx.service.TableMessage().Name()).
					AsRequired()
			})

			message.CreateField("update_mask", func(field *presentation.Field) {
				field.SetType("google.protobuf.FieldMask").SetTopComment("The list of fields to update.")
			})
		},
	)

	return nil
}

var autofilledTimestampColumnNames = map[string]bool{
	"deleted_at": true, "delete_time": true,
	"updated_at": true, "update_time": true,
	"created_at": true, "create_time": true,
}

func (e *Exporter) columnAutofilled(col *schema.Column) bool {
	if col.IsAutoincrement {
		return true
	}

	if col.Type.IsDatetime || col.Type.IsDate {
		if autofilledTimestampColumnNames[col.Name.Lower().Value] {
			return true
		}
	}

	if col.DefaultRaw.Valid && col.IsPrimaryKey() {
		return true
	}

	return false
}
