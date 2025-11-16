package pg

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"

	"github.com/artarts36/gds"

	"github.com/artarts36/db-exporter/internal/infrastructure/conn"
)

type partitionRel struct {
	ChildTable  gds.String `db:"child_table"`
	ParentTable gds.String `db:"parent_table"`
}

func (l *Loader) loadPartitions(ctx context.Context, cn *conn.Connection, tableNames []string) ([]*partitionRel, error) {
	query := `SELECT
    inhrelid::regclass AS child_table,
    inhparent::regclass AS parent_table
FROM pg_inherits
    JOIN pg_class parent            ON pg_inherits.inhparent = parent.oid
    JOIN pg_class child             ON pg_inherits.inhrelid  = child.oid
WHERE
    child.relpartbound IS NOT NULL AND 
	inhrelid::regclass IN (?)
ORDER BY inhrelid
`
	db, err := cn.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	var partitions []*partitionRel

	query, args, err := sqlx.In(query, tableNames)
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	query = db.Rebind(query)
	err = db.SelectContext(ctx, &partitions, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}

	return partitions, nil
}
