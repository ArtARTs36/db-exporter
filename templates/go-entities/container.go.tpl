package {{ package.Name }}

import (
    "github.com/jmoiron/sqlx"
)

type {{ containerName }} struct {
{% for repo in schema.Repositories %}    {{ repo.Name }} *{{ repo.Name }}{% if loop.last == false %}
{% endif %}{% endfor %}
}

func New{{ containerName }}(db *sqlx.DB) *{{ containerName }} {
    return &{{ containerName }}{
{% for repo in schema.Repositories %}        {{ repo.Name }}: New{{ repo.Name }}(db),{% if loop.last == false %}
{% endif %}{% endfor %}
    }
}
