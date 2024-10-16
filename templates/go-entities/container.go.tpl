package {{ _file.Package.Name }}

import (
	"github.com/jmoiron/sqlx"
)

type {{ containerName }} struct {
{% for repo in schema.Repositories %}	{{ repo.Name }} {{ spaces_after(repo.Name, schema.RepoNameMaxLength) }}*{{ repo.Name }}{% if loop.last == false %}
{% endif %}{% endfor %}
}

func New{{ containerName }}(db *sqlx.DB) *{{ containerName }} {
	return &{{ containerName }}{
{% for repo in schema.Repositories %}		{{ repo.Name }}: {{ spaces_after(repo.Name, schema.RepoNameMaxLength) }}New{{ repo.Name }}(db),{% if loop.last == false %}
{% endif %}{% endfor %}
	}
}
