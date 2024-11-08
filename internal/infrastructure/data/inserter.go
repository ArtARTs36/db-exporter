package data

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/infrastructure/conn"

	"github.com/doug-martin/goqu/v9"

	"github.com/artarts36/db-exporter/internal/schema"
)

type Inserter struct {
}

func NewInserter() *Inserter {
	return &Inserter{}
}

func (i *Inserter) Insert(
	ctx context.Context,
	cn *conn.Connection,
	table string,
	dataset []map[string]interface{},
) (int64, error) {
	extCtx, err := cn.ExtContext(ctx)
	if err != nil {
		return 0, err
	}

	rows := make([]interface{}, 0, len(dataset))
	for _, row := range dataset {
		rows = append(rows, row)
	}

	q, _, err := goqu.Insert(table).Rows(rows...).ToSQL()
	if err != nil {
		return 0, fmt.Errorf("failed to build insert query: %w", err)
	}

	res, err := extCtx.ExecContext(ctx, q)
	if err != nil {
		return 0, fmt.Errorf("failed to insert dataset into database: %w", err)
	}

	return res.RowsAffected()
}

func (i *Inserter) Upsert(
	ctx context.Context,
	cn *conn.Connection,
	table *schema.Table,
	dataset []map[string]interface{},
) (int64, error) {
	extCtx, err := cn.ExtContext(ctx)
	if err != nil {
		return 0, err
	}

	rows := make([]interface{}, 0, len(dataset))
	for _, row := range dataset {
		rows = append(rows, row)
	}

	updateRecord := goqu.Record{}
	for col := range dataset[0] {
		updateRecord[col] = goqu.I(fmt.Sprintf("excluded.%s", col))
	}

	q, _, err := goqu.
		Insert(table.Name.Value).
		Rows(rows...).
		OnConflict(goqu.DoUpdate(
			table.PrimaryKey.ColumnsNames.Join(",").Value,
			updateRecord,
		)).
		ToSQL()
	if err != nil {
		return 0, fmt.Errorf("failed to build insert query: %w", err)
	}

	res, err := extCtx.ExecContext(ctx, q)
	if err != nil {
		return 0, fmt.Errorf("failed to insert dataset into database: %w", err)
	}

	return res.RowsAffected()
}
