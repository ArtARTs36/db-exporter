package {{ package.Name }}

import (
    "context"
    "fmt"

    "github.com/doug-martin/goqu/v9"
    "github.com/jmoiron/sqlx"
)

const (
{% for repo in schema.Repositories %}    table{{ repo.Entity.Table.Name.Pascal().Value }} = "{{ repo.Entity.Table.Name.Value }}"{% if loop.last == false %}
{% endif %}{% endfor %}
)

{% for repo in schema.Repositories %}type {{ repo.Name }} struct {
    db *sqlx.DB
}

func New{{ repo.Name }}(db *sqlx.DB) *{{ repo.Name }} {
    return &{{ repo.Name }}{db: db}
}

func (repo *{{ repo.Name }}) Create(
    ctx context.Context,
    {{ repo.Entity.Table.Name.Singular().Camel().Value }} *{{ repo.EntityCall }},
) (*{{ repo.EntityCall }}, error) {
    query, _, err := goqu.Insert(table{{ repo.Entity.Table.Name.Pascal().Value }}).Rows({{ repo.Entity.Table.Name.Singular().Camel().Value }}).Returning("*").ToSQL()
    if err != nil {
        return nil, fmt.Errorf("failed to build insert query: %w", err)
    }

    var created {{ repo.EntityCall }}

    err = repo.db.GetContext(ctx, &created, query)
    if err != nil {
        return nil, fmt.Errorf("failed to execute query: %w", err)
    }

    return &created, nil
}
{% endfor %}
