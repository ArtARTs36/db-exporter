package schema

import "database/sql"

type Column struct {
	Name       String         `db:"name"`
	TableName  String         `db:"table_name"`
	Type       String         `db:"type"`
	Nullable   bool           `db:"nullable"`
	PrimaryKey sql.NullString `db:"-"`
	ForeignKey *ForeignKey    `db:"-"`
	Comment    String         `db:"comment"`
}

func (c *Column) IsPrimaryKey() bool {
	return c.PrimaryKey.Valid
}

func (c *Column) HasForeignKey() bool {
	return c.ForeignKey != nil
}
