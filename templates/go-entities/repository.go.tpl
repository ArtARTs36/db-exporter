package {{ package.Name }}

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
{% if not entityPackage.IsCurrent(package) %}
	"{{ entityPackage.FullName }}"{% endif %}
)

const (
{% for repo in schema.Repositories %}	table{{ repo.Entity.Table.Name.Pascal().Value }} = "{{ repo.Entity.Table.Name.Value }}"{% if loop.last == false %}
{% endif %}{% endfor %}
)

{% for repo in schema.Repositories %}type {{ repo.Name }} struct {
	db *sqlx.DB
}{% endfor %}

{% for repo in schema.Repositories %}type {{ repo.Filters.List.Name }} struct {
{% for prop in repo.Filters.List.Properties.List %}	{{ prop.PluralName }} {{ spaces_after(prop.Name, entity.Properties.MaxPropNameLength) }}[]{{ prop.Type }}{% if loop.last == false %}
{% endif %}{% endfor %}
}{% endfor %}

{% for repo in schema.Repositories %}func New{{ repo.Name }}(db *sqlx.DB) *{{ repo.Name }} {
	return &{{ repo.Name }}{db: db}
}

func (repo *{{ repo.Name }}) List(
    ctx context.Context,
    filter *{{ repo.Filters.List.Name }},
) ([]*{{ repo.EntityCall }}, error) {
    var ents []*{{ repo.EntityCall }}

    query := goqu.From(table{{ repo.Entity.Table.Name.Pascal().Value }}).Select()

{% if repo.Filters.List.Properties.List | length > 0 %}    if filter != nil {
{% for prop in repo.Filters.List.Properties.List %}        if len(filter.{{ prop.PluralName }}) > 0 {
            query = query.Where(goqu.C("{{ prop.Column.Name.Value }}").In(filter.{{ prop.PluralName }}))
        }{% if loop.last == false %}
{% endif %}{% endfor %}
    }{% endif %}

    q, args, err := query.ToSQL()
    if err != nil {
        return nil, fmt.Errorf("failed to build query: %w", err)
    }

    err = repo.db.SelectContext(ctx, &ents, q, args...)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return []*{{ repo.EntityCall }}{}, nil
        }
        return nil, fmt.Errorf("failed to execute query: %w", err)
    }

    return ents, nil
}

func (repo *{{ repo.Name }}) Create(
	ctx context.Context,
	{{ repo.Entity.AsVarName }} *{{ repo.EntityCall }},
) (*{{ repo.EntityCall }}, error) {
	query, _, err := goqu.Insert(table{{ repo.Entity.Table.Name.Pascal().Value }}).Rows({{ repo.Entity.AsVarName }}).Returning("*").ToSQL()
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

func (repo *{{ repo.Name }}) Update(
	ctx context.Context,
	{{ repo.Entity.AsVarName }} *{{ repo.EntityCall }},
) (*{{ repo.EntityCall }}, error) {
    query, _, err := goqu.Update(table{{ repo.Entity.Table.Name.Pascal().Value }}).
        Set({{ repo.Entity.AsVarName }}).
        Returning("*").
        ToSQL()
    if err != nil {
        return nil, fmt.Errorf("failed to build update query: %w", err)
    }

    var updated {{ repo.EntityCall }}

    err = repo.db.GetContext(ctx, &updated, query)
    if err != nil {
        return nil, fmt.Errorf("failed to execute query: %w", err)
    }

    return &updated, nil
}{% endfor %}
