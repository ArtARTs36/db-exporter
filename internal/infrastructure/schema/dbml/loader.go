package dbml

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/artarts36/db-exporter/internal/infrastructure/conn"
	"github.com/artarts36/db-exporter/internal/infrastructure/sqltype"
	"github.com/artarts36/dbml-go/core"
	"github.com/artarts36/dbml-go/parser"
	"github.com/artarts36/gds"
	"os"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/dbml"
)

type Loader struct {
}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) Load(ctx context.Context, cn *conn.Connection) (*schema.Schema, error) {
	f, err := os.OpenFile(cn.Database().DSN.Value, os.O_RDONLY, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", cn.Database().DSN.Value, err)
	}

	parsedDBML, err := parser.Parse(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dbml: %w", err)
	}

	return l.buildSchema(parsedDBML)
}

func (l *Loader) buildSchema(parsedDBML *core.DBML) (*schema.Schema, error) {
	sch := schema.NewSchema(schema.DatabaseDriverDBML)
	sch.Enums = l.collectEnums(parsedDBML)

	for _, tbl := range parsedDBML.Tables {
		table := &schema.Table{
			Name:           *gds.NewString(tbl.Name),
			Comment:        tbl.Note,
			ForeignKeys:    map[string]*schema.ForeignKey{},
			UsingSequences: map[string]*schema.Sequence{},
			UniqueKeys:     map[string]*schema.UniqueKey{},
		}

		sch.Tables.Add(table)

		columns := make([]*schema.Column, 0, len(tbl.Columns))

		for _, col := range tbl.Columns {
			column := &schema.Column{
				Name:            *gds.NewString(col.Name),
				DataType:        sqltype.MapDBMLType(col.Type),
				TableName:       table.Name,
				Nullable:        col.Settings.Null,
				IsAutoincrement: col.Settings.Increment,
				DefaultRaw: sql.NullString{
					Valid:  col.Settings.Default.Raw != "",
					String: col.Settings.Default.Raw,
				},
				Default:        l.parseDefaultValue(col.Settings.Default),
				UsingSequences: map[string]*schema.Sequence{},
				Comment:        *gds.NewString(col.Settings.Note),
			}

			if col.Settings.Unique {
				uk := &schema.UniqueKey{
					Name:         *gds.NewString(fmt.Sprintf("%s_%s_uk", table.Name.Value, col.Name)),
					ColumnsNames: gds.NewStrings(col.Name),
				}

				column.UniqueKey, table.UniqueKeys[uk.Name.Value] = uk, uk
			}

			if col.Settings.PK {
				pk := &schema.PrimaryKey{
					Name:         *gds.NewString(fmt.Sprintf("%s_%s_pk", table.Name.Value, col.Name)),
					ColumnsNames: gds.NewStrings(col.Name),
				}

				column.PrimaryKey, table.PrimaryKey = pk, pk
			}

			if col.Settings.Increment {
				incSeq := schema.CreateSequenceForColumn(column)
				incSeq.Inc()
				column.UsingSequences[incSeq.Name] = incSeq
				table.UsingSequences[incSeq.Name] = incSeq
				sch.Sequences[incSeq.Name] = incSeq
			}

			columns = append(columns, column)
		}

		table.Columns = columns
	}

	for i, ref := range parsedDBML.Refs {
		if len(ref.Relationships) != 1 {
			return nil, fmt.Errorf("ref[%d] must have one relationship", i)
		}

		relation := ref.Relationships[0]

		from := relation.From
		to := relation.To

		switch relation.Type {
		case core.OneToMany, core.OneToOne:
			from = relation.To
			to = relation.From
		case core.ManyToOne, core.None:
			from = relation.From
			to = relation.To
		}

		fromSubject, err := l.getRelationSubject(sch, from)
		if err != nil {
			return nil, fmt.Errorf("failed to get relation from ref[%d]: %w", i, err)
		}

		toSubject, err := l.getRelationSubject(sch, to)
		if err != nil {
			return nil, fmt.Errorf("failed to get relation from ref[%d]: %w", i, err)
		}

		fk := l.buildForeignKey(fromSubject, toSubject)
		fromSubject.Column.ForeignKey = fk
		fromSubject.Table.ForeignKeys[fk.Name.Value] = fk
	}

	return sch, nil
}

func (l *Loader) parseDefaultValue(raw core.ColumnDefault) *schema.ColumnDefault {
	if raw.Raw == "" {
		return nil
	}

	switch raw.Type { //nolint:exhaustive // not need
	case core.ColumnDefaultTypeString:
		return &schema.ColumnDefault{
			Type:  schema.ColumnDefaultTypeValue,
			Value: raw.Value,
		}
	case core.ColumnDefaultTypeNumber:
		return &schema.ColumnDefault{
			Type:  schema.ColumnDefaultTypeValue,
			Value: raw.Value,
		}
	case core.ColumnDefaultTypeExpression:
		return &schema.ColumnDefault{
			Type:  schema.ColumnDefaultTypeFunc,
			Value: raw.Value,
		}
	}

	return nil
}

func (l *Loader) collectEnums(parsedDBML *core.DBML) map[string]*schema.Enum {
	enums := map[string]*schema.Enum{}

	for _, en := range parsedDBML.Enums {
		enum := &schema.Enum{
			Name:          gds.NewString(en.Name),
			Values:        make([]string, 0, len(en.Values)),
			UsingInTables: make([]string, 0),
		}

		enums[en.Name] = enum

		for _, value := range en.Values {
			enum.Values = append(enum.Values, value.Name)
		}
	}

	return enums
}

func (l *Loader) buildForeignKey(from *dbmlRelationSubject, to *dbmlRelationSubject) *schema.ForeignKey {
	return &schema.ForeignKey{
		Name: *from.Table.Name.Append("_").
			Append(from.Column.Name.Value).
			Append("_").
			Append(to.Table.Name.Value).
			Append("_").
			Append(to.Column.Name.Value).
			Append("_fk"),

		Table:         from.Table.Name,
		ColumnsNames:  gds.NewStrings(from.Column.Name.Value),
		ForeignTable:  to.Table.Name,
		ForeignColumn: to.Column.Name,
	}
}

type dbmlRelationSubject struct {
	Table  *schema.Table
	Column *schema.Column
}

func (l *Loader) getRelationSubject(sch *schema.Schema, subj string) (*dbmlRelationSubject, error) {
	rel, err := dbml.ParseRelationSubject(subj)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ref: %w", err)
	}

	table, ok := sch.Tables.Get(*gds.NewString(rel.Table))
	if !ok {
		return nil, fmt.Errorf("table %q not found", rel.Table)
	}

	column := table.GetColumn(rel.Column)
	if column == nil {
		return nil, fmt.Errorf("column %q not found in table %q", table.Name, rel.Column)
	}

	return &dbmlRelationSubject{
		Table:  table,
		Column: column,
	}, nil
}
