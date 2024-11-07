package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/artarts36/dbml-go/core"
	"github.com/artarts36/dbml-go/parser"
	"github.com/artarts36/dbml-go/scanner"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/dbml"
	"github.com/artarts36/db-exporter/internal/shared/ds"
)

type DBMLLoader struct { //nolint:revive // 'DB' part of the name
}

func NewDBMLLoader() *DBMLLoader {
	return &DBMLLoader{}
}

func (l *DBMLLoader) Load(_ context.Context, conn *Connection) (*schema.Schema, error) {
	f, err := os.OpenFile(conn.cfg.DSN, os.O_RDONLY, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", conn.cfg.DSN, err)
	}

	parsedDBML, err := parser.NewParser(scanner.NewScanner(f)).Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse dbml: %w", err)
	}

	return l.buildSchema(parsedDBML)
}

func (l *DBMLLoader) buildSchema(parsedDBML *core.DBML) (*schema.Schema, error) {
	sch := schema.NewSchema()
	sch.Enums = l.collectEnums(parsedDBML)

	for _, tbl := range parsedDBML.Tables {
		table := &schema.Table{
			Name:           *ds.NewString(tbl.Name),
			Comment:        tbl.Note,
			ForeignKeys:    map[string]*schema.ForeignKey{},
			UsingSequences: map[string]*schema.Sequence{},
		}

		sch.Tables.Add(table)

		columns := make([]*schema.Column, 0, len(tbl.Columns))

		for _, col := range tbl.Columns {
			column := &schema.Column{
				Name:            *ds.NewString(col.Name),
				Type:            *ds.NewString(col.Type),
				PreparedType:    l.mapGoType(col.Type),
				TableName:       table.Name,
				Nullable:        col.Settings.Null,
				IsAutoincrement: col.Settings.Increment,
				DefaultRaw: sql.NullString{
					Valid:  col.Settings.Default != "",
					String: col.Settings.Default,
				},
				UsingSequences: map[string]*schema.Sequence{},
				Comment:        *ds.NewString(col.Settings.Note),
			}

			if col.Settings.Unique {
				uk := &schema.UniqueKey{
					Name:         *ds.NewString(fmt.Sprintf("%s_%s_uk", table.Name.Value, col.Name)),
					ColumnsNames: ds.NewStrings(col.Name),
				}

				column.UniqueKey, table.UniqueKeys[uk.Name.Value] = uk, uk
			}

			if col.Settings.PK {
				pk := &schema.PrimaryKey{
					Name:         *ds.NewString(fmt.Sprintf("%s_%s_pk", table.Name.Value, col.Name)),
					ColumnsNames: ds.NewStrings(col.Name),
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

func (l *DBMLLoader) collectEnums(parsedDBML *core.DBML) map[string]*schema.Enum {
	enums := map[string]*schema.Enum{}

	for _, en := range parsedDBML.Enums {
		enum := &schema.Enum{
			Name:   ds.NewString(en.Name),
			Values: make([]string, 0, len(en.Values)),
		}

		enums[en.Name] = enum

		for _, value := range en.Values {
			enum.Values = append(enum.Values, value.Name)
		}
	}

	return enums
}

func (l *DBMLLoader) buildForeignKey(from *dbmlRelationSubject, to *dbmlRelationSubject) *schema.ForeignKey {
	return &schema.ForeignKey{
		Name: *from.Table.Name.Append("_").
			Append(from.Column.Name.Value).
			Append("_").
			Append(to.Table.Name.Value).
			Append("_").
			Append(to.Column.Name.Value).
			Append("_fk"),

		Table:         from.Table.Name,
		ColumnsNames:  ds.NewStrings(from.Column.Name.Value),
		ForeignTable:  to.Table.Name,
		ForeignColumn: to.Column.Name,
	}
}

type dbmlRelationSubject struct {
	Table  *schema.Table
	Column *schema.Column
}

func (l *DBMLLoader) getRelationSubject(sch *schema.Schema, subj string) (*dbmlRelationSubject, error) {
	rel, err := dbml.ParseRelationSubject(subj)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ref: %w", err)
	}

	table, ok := sch.Tables.Get(*ds.NewString(rel.Table))
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

func (l *DBMLLoader) mapGoType(rawType string) schema.DataType {
	switch rawType {
	case "integer":
		return schema.DataTypeInteger
	case "varchar", "text":
		return schema.DataTypeString
	case "bool", "boolean":
		return schema.DataTypeBoolean
	}
	return schema.DataTypeString
}
