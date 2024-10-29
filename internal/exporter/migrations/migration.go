package migrations

import "github.com/artarts36/db-exporter/internal/shared/ds"

type MigrationMaker interface {
	MakeSingle(index int, tableName ds.String) *MigrationMeta
	MakeMultiple() *MigrationMeta
}

type (
	MakeSingleMigrationFunc   func(index int, tableName ds.String) *MigrationMeta
	MakeMultipleMigrationFunc func() *MigrationMeta
)

type MigrationMeta struct {
	Filename string
	Attrs    map[string]interface{}
}

type Migration struct {
	Meta map[string]interface{}

	UpQueries   []string
	DownQueries []string
}

type FuncMigrationMaker struct {
	makeSingle   MakeSingleMigrationFunc
	makeMultiple MakeMultipleMigrationFunc
}

func NewFuncMigrationMaker(single MakeSingleMigrationFunc, multiple MakeMultipleMigrationFunc) *FuncMigrationMaker {
	return &FuncMigrationMaker{
		makeSingle:   single,
		makeMultiple: multiple,
	}
}

func (m *FuncMigrationMaker) MakeSingle(index int, tableName ds.String) *MigrationMeta {
	return m.makeSingle(index, tableName)
}

func (m *FuncMigrationMaker) MakeMultiple() *MigrationMeta {
	return m.makeMultiple()
}
