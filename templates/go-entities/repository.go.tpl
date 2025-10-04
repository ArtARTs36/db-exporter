{% if schema.WithMocks %}{% include '@embed/go-entities/gomock_call.go.tpl' with {'_file': _file, 'repositories': schema.Repositories} only %}
{% endif %}package {{ _file.Package.Name }}

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
){% if schema.GenInterfaces %}

{% include '@embed/go-entities/repository_interfaces.go.tpl' with {'repositories': schema.Repositories, '_file': _file} only %}
{% else %}
{% endif %}
{% for repo in schema.Repositories %}type {{ repo.Name }} struct {
	db *sqlx.DB
}{% endfor %}
{% if schema.GenFilters %}
{% for repo in schema.Repositories %}{% include '@embed/go-entities/repository_filters.go.tpl' with {'filters': repo.Filters} only %}{% endfor %}
{% endif %}
{% for repo in schema.Repositories %}func New{{ repo.Name }}(db *sqlx.DB) *{{ repo.Name }} {
	return &{{ repo.Name }}{db: db}
}

func (repo *{{ repo.Name }}) Get(
	ctx context.Context,
	filter *{{ repo.Filters.Get.Call(_file.Package) }},
) (*{{ repo.Entity.Call(_file.Package) }}, error) {
	var ent {{ repo.Entity.Call(_file.Package) }}

	query := goqu.From(table{{ repo.Entity.Table.Name.Pascal().Value }}).Select().Limit(1)

{% if repo.Filters.Get.Properties.List | length > 0 %}	if filter != nil {
{% for prop in repo.Filters.Get.Properties.List %}		{% if prop.IsString() %}if len(filter.{{ prop.Name }}) > 0 {{ figure_brace_opened }}{% else %}if filter.{{ prop.Name }} > 0 {{ figure_brace_opened }}{% endif %}
			query = query.Where(goqu.C("{{ prop.Column.Name.Value }}").Eq(filter.{{ prop.Name }}))
		}{% if loop.last == false %}
{% endif %}{% endfor %}
	}{% endif %}

	q, args, err := query.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = repo.db.GetContext(ctx, &ent, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &ent, nil
}

func (repo *{{ repo.Name }}) List(
	ctx context.Context,
	filter *{{ repo.Filters.List.Call(_file.Package) }},
) ([]*{{ repo.Entity.Call(_file.Package) }}, error) {
	var ents []*{{ repo.Entity.Call(_file.Package) }}

	query := goqu.From(table{{ repo.Entity.Table.Name.Pascal().Value }}).Select()

{% if repo.Filters.List.Properties.List | length > 0 %}	if filter != nil {
{% for prop in repo.Filters.List.Properties.List %}		if len(filter.{{ prop.PluralName }}) > 0 {
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
			return []*{{ repo.Entity.Call(_file.Package) }}{}, nil
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return ents, nil
}

func (repo *{{ repo.Name }}) Create(
	ctx context.Context,
	{{ repo.Entity.AsVarName }} *{{ repo.Entity.Call(_file.Package) }},
) (*{{ repo.Entity.Call(_file.Package) }}, error) {
	query, _, err := goqu.Insert(table{{ repo.Entity.Table.Name.Pascal().Value }}).Rows({{ repo.Entity.AsVarName }}).Returning("*").ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build insert query: %w", err)
	}

	var created {{ repo.Entity.Call(_file.Package) }}

	err = repo.db.GetContext(ctx, &created, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &created, nil
}

func (repo *{{ repo.Name }}) Update(
	ctx context.Context,
	{{ repo.Entity.AsVarName }} *{{ repo.Entity.Call(_file.Package) }},
) (*{{ repo.Entity.Call(_file.Package) }}, error) {
	query, _, err := goqu.Update(table{{ repo.Entity.Table.Name.Pascal().Value }}).
		Set({{ repo.Entity.AsVarName }}).
		Returning("*").
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build update query: %w", err)
	}

	var updated {{ repo.Entity.Call(_file.Package) }}

	err = repo.db.GetContext(ctx, &updated, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &updated, nil
}

func (repo *{{ repo.Name }}) Delete(
	ctx context.Context,
	filter *{{ repo.Filters.Delete.Call(_file.Package) }},
) (count int64, err error) {
	query := goqu.From(table{{ repo.Entity.Table.Name.Pascal().Value }}).Delete()

{% if repo.Filters.Get.Properties.List | length > 0 %}	if filter != nil {
{% for prop in repo.Filters.Get.Properties.List %}		if len(filter.{{ prop.PluralName }}) > 0 {
			query = query.Where(goqu.C("{{ prop.Column.Name.Value }}").In(filter.{{ prop.PluralName }}))
		}{% if loop.last == false %}
{% endif %}{% endfor %}
	}{% endif %}

	q, args, err := query.ToSQL()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	res, err := repo.db.ExecContext(ctx, q, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}
	count, err = res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get affected rows: %w", err)
	}

	return
}{% endfor %}
