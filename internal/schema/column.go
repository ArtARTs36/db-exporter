package schema

type Column struct {
	Name         String      `db:"name"`
	TableName    String      `db:"table_name"`
	Type         String      `db:"type"`
	Nullable     bool        `db:"nullable"`
	PrimaryKey   *PrimaryKey `db:"-"`
	UniqueKey    *UniqueKey  `db:"-"`
	ForeignKey   *ForeignKey `db:"-"`
	Comment      String      `db:"comment"`
	PreparedType ColumnType  `db:"-"`
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
