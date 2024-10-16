package {{ _file.Package.Name }}

import (
	"github.com/jmoiron/sqlx"
)

type {{ schema.Container.Name }} struct {
{% for repo in schema.Repositories %}	{{ repo.Interface.Name }} {{ spaces_after(repo.Interface.Name, schema.RepoInterfaceNameMaxLength) }}{{ repo.Interface.Name }}{% if loop.last == false %}
{% endif %}{% endfor %}
}

func New{{ schema.Container.Name }}(db *sqlx.DB) *{{ schema.Container.Name }} {
	return &{{ schema.Container.Name }}{
{% for repo in schema.Repositories %}		{{ repo.Interface.Name }}: {{ spaces_after(repo.Name, schema.RepoNameMaxLength) }}New{{ repo.Name }}(db),{% if loop.last == false %}
{% endif %}{% endfor %}
	}
}
