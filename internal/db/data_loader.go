package db

import (
	"context"
	"fmt"
	"log/slog"
)

type TableData []map[string]interface{}

type DataLoader struct {
	db *Connection
}

func NewDataLoader(conn *Connection) *DataLoader {
	return &DataLoader{db: conn}
}

func (l *DataLoader) Load(ctx context.Context, table string) (TableData, error) {
	data := make(TableData, 0)

	q := fmt.Sprintf("select * from %s", table)

	db, err := l.db.Connect(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", rows.Err())
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("failed to close rows: %s", closeErr))
		}
	}()

	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch columns: %w", err)
	}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if scanErr := rows.Scan(columnPointers...); scanErr != nil {
			return nil, fmt.Errorf("failed to scan: %w", scanErr)
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val, _ := columnPointers[i].(*interface{})

			m[colName] = *val
		}

		data = append(data, m)
	}

	return data, nil
}
