package schema

import (
	"database/sql"
	"github.com/artarts36/gds"
)

type Column struct {
	Name           gds.String           `db:"name"`
	TableName      gds.String           `db:"table_name"`
	Type           gds.String           `db:"type"`
	Nullable       bool                 `db:"nullable"`
	PrimaryKey     *PrimaryKey          `db:"-"`
	UniqueKey      *UniqueKey           `db:"-"`
	ForeignKey     *ForeignKey          `db:"-"`
	Comment        gds.String           `db:"comment"`
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
