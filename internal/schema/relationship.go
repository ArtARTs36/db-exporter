package schema

import "fmt"

type RelationshipType string

const (
	RelationshipTypeOneToOne   RelationshipType = "one-to-one"
	RelationshipTypeManyToMany RelationshipType = "many-to-many"
	RelationShipOneToMany      RelationshipType = "one-to-many"
)

type Relationship struct {
	Type RelationshipType

	OwnerTable  *Table
	OwnerColumn *Column

	RelatedTable  *Table
	RelatedColumn *Column
}

// DetermineRelationships
// Returns map[owner table][]RelationShip
func DetermineRelationships(schema *Schema) (map[string][]*Relationship, error) {
	relationships := make(map[string][]*Relationship)

	allocate := func(ownerTableName string) {
		if _, ok := relationships[ownerTableName]; !ok {
			relationships[ownerTableName] = make([]*Relationship, 0)
		}
	}

	for _, table := range schema.Tables.List() {
		if len(table.ForeignKeys) == 0 {
			continue
		}

		// naive search many-to-many
		if len(table.ForeignKeys) == 2 && len(table.Columns) == 2 {
			firstTable, ok := schema.Tables.Get(table.Columns[0].ForeignKey.ForeignTable)
			if !ok {
				return nil, fmt.Errorf("table %s not found", table.Columns[0].ForeignKey.ForeignTable)
			}

			secondTable, ok := schema.Tables.Get(table.Columns[1].ForeignKey.ForeignTable)
			if !ok {
				return nil, fmt.Errorf("table %s not found", table.Columns[1].ForeignKey.ForeignTable)
			}

			allocate(firstTable.Name.Value)
			relationships[firstTable.Name.Value] = append(relationships[firstTable.Name.Value], &Relationship{
				Type: RelationshipTypeManyToMany,

				OwnerTable:  firstTable,
				OwnerColumn: firstTable.GetColumn(table.Columns[0].ForeignKey.ForeignColumn.Value),

				RelatedTable:  secondTable,
				RelatedColumn: secondTable.GetColumn(table.Columns[1].ForeignKey.ForeignColumn.Value),
			})

			allocate(secondTable.Name.Value)
			relationships[secondTable.Name.Value] = append(relationships[secondTable.Name.Value], &Relationship{
				Type: RelationshipTypeManyToMany,

				OwnerTable:  secondTable,
				OwnerColumn: secondTable.GetColumn(table.Columns[1].ForeignKey.ForeignColumn.Value),

				RelatedTable:  firstTable,
				RelatedColumn: firstTable.GetColumn(table.Columns[0].ForeignKey.ForeignColumn.Value),
			})

			continue
		}

		relationships[table.Name.Value] = make([]*Relationship, 0)

		for _, foreignKey := range table.ForeignKeys {
			relatedTable, ok := schema.Tables.Get(foreignKey.ForeignTable)
			if !ok {
				return nil, fmt.Errorf("table %s not found", foreignKey.Table)
			}

			ship := &Relationship{
				Type: RelationshipTypeOneToOne,

				OwnerTable:  table,
				OwnerColumn: table.GetColumn(foreignKey.ColumnsNames.First()),

				RelatedTable:  relatedTable,
				RelatedColumn: relatedTable.GetColumn(foreignKey.ForeignColumn.Value),
			}

			relationships[table.Name.Value] = append(relationships[table.Name.Value], ship)

			allocate(foreignKey.ForeignTable.Value)
			relationships[relatedTable.Name.Value] = append(relationships[foreignKey.ForeignTable.Value], &Relationship{
				Type: RelationShipOneToMany,

				OwnerTable:  relatedTable,
				OwnerColumn: relatedTable.GetColumn(foreignKey.ForeignColumn.Value),

				RelatedTable:  table,
				RelatedColumn: table.GetColumn(foreignKey.ColumnsNames.First()),
			})
		}
	}

	return relationships, nil
}
