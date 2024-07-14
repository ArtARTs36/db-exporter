package db

import (
	"context"
	"fmt"
	"github.com/doug-martin/goqu/v9"
)

type Inserter struct {
	db *Connection
}

func NewInserter(db *Connection) *Inserter {
	return &Inserter{db: db}
}

func (i *Inserter) Insert(ctx context.Context, table string, dataset []map[string]interface{}) error {
	db, err := i.db.Connect(ctx)
	if err != nil {
		return err
	}

	rows := make([]interface{}, 0, len(dataset))
	for _, row := range dataset {
		fmt.Println(row)
		rows = append(rows, row)
	}

	q, _, err := goqu.Insert(table).Rows(rows...).ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	_, err = db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("failed to insert dataset into database: %w", err)
	}

	return nil
}
