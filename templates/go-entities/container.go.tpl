package {{ package.Name }}

import (
    "github.com/jmoiron/sqlx"
)

type {{ container_name }} struct {
{% for repo in schema.Repositories %}    {{ repo.Name }} *{{ repo.Name }}{% if loop.last == false %}
{% endif %}{% endfor %}
}

func New{{ container_name }}(db *sqlx.DB) *Group {
    return &Group{
{% for repo in schema.Repositories %}        {{ repo.Name }}: New{{ repo.Name }}(db),{% if loop.last == false %}
{% endif %}{% endfor %}
    }
}
