package schema

import (
	"database/sql"
	"github.com/artarts36/db-exporter/internal/shared/ds"
)

type Column struct {
	Name           ds.String            `db:"name"`
	TableName      ds.String            `db:"table_name"`
	Type           ds.String            `db:"type"`
	Nullable       bool                 `db:"nullable"`
	PrimaryKey     *PrimaryKey          `db:"-"`
	UniqueKey      *UniqueKey           `db:"-"`
	ForeignKey     *ForeignKey          `db:"-"`
	Comment        ds.String            `db:"comment"`
	PreparedType   DataType             `db:"-"`
	DefaultRaw     sql.NullString       `db:"default_value"`
	Default        *ColumnDefault       `db:"-"`
	UsingSequences map[string]*Sequence `db:"-"`
	Enum           *Enum                `db:"-"`

	IsAutoincrement bool `db:"-"`
}

func (c *Column) IsPrimaryKey() bool {
	return c.PrimaryKey != nil
}

func (c *Column) IsUniqueKey() bool {
	return c.UniqueKey != nil
}

func (c *Column) HasForeignKey() bool {
	return c.ForeignKey != nil
}

func (c *Column) IsUniqueOrPrimaryKey() bool {
	return c.IsUniqueKey() || c.IsPrimaryKey()
}
