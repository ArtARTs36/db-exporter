package pg

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/infrastructure/conn"
	"github.com/artarts36/db-exporter/internal/infrastructure/sqltype"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/gds"
)

type pgDomain struct {
	Name           string `db:"name"`
	DataType       string `db:"data_type"`
	ConstraintName string `db:"constraint_name"`
	CheckClause    string `db:"check_clause"`
}

func (l *Loader) loadDomains(ctx context.Context, cn *conn.Connection) (*gds.Map[string, *schema.Domain], error) {
	query := `
SELECT
    d.domain_name as name,
    d.data_type,
    dc.constraint_name as constraint_name,
    pg_get_constraintdef(c.oid) as check_clause
FROM information_schema.domains d
JOIN
    information_schema.domain_constraints dc ON dc.domain_name = d.domain_name
JOIN
    pg_constraint c ON c.conname = dc.constraint_name
WHERE d.domain_schema = $1 AND dc.constraint_schema = $1`

	db, err := cn.Connect(ctx)
	if err != nil {
		return nil, err
	}

	domains := make([]*pgDomain, 0)

	err = db.SelectContext(ctx, &domains, query, cn.Database().Schema)
	if err != nil {
		return nil, fmt.Errorf("select domains: %w", err)
	}

	schemaDomains := gds.NewMap[string, *schema.Domain]()

	for _, d := range domains {
		schemaDomains.Set(d.Name, &schema.Domain{
			Name:           d.Name,
			DataType:       sqltype.MapPGType(d.DataType),
			ConstraintName: d.ConstraintName,
			CheckClause:    d.CheckClause,
		})
	}

	return schemaDomains, nil
}
