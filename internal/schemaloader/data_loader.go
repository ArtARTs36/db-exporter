package schemaloader

import (
	"context"
	"fmt"
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

	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch columns: %w", err)
	}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		data = append(data, m)
	}

	return data, nil
}
