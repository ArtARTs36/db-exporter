package schema

import (
	"database/sql"
	"github.com/artarts36/gds"
)

type Column struct {
	Name            gds.String           `db:"name"`       // Name of column.
	TableName       gds.String           `db:"table_name"` // Name of table.
	DataType        DataType             `db:"-"`          // Column datatype
	TypeRaw         gds.String           `db:"type_raw"`
	CharacterLength int16                `db:"character_length"` // Length of character.
	Nullable        bool                 `db:"nullable"`         // True, if column supports null value
	PrimaryKey      *PrimaryKey          `db:"-"`                // Primary key.
	UniqueKey       *UniqueKey           `db:"-"`                // Unique key.
	ForeignKey      *ForeignKey          `db:"-"`                // The foreign key referenced by this column.
	Comment         gds.String           `db:"comment"`          // Comment for column.
	DefaultRaw      sql.NullString       `db:"default_value"`
	Default         *ColumnDefault       `db:"-"`
	UsingSequences  map[string]*Sequence `db:"-"`           // Map of using sequences.
	DomainName      sql.NullString       `db:"domain_name"` // Name of domain (for PostgreSQL).

	Enum   *Enum   `db:"-"` // List of values this column.
	Domain *Domain `db:"-"` // Domain  (for PostgreSQL)

	// True, if the column supports auto-increment.
	// Postgres sequence is also auto-increment=true.
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
