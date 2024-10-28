{% include 'go-entities/go_file_header.go.tpl' with {'_file': _file, 'sharedImports': ['github.com/jmoiron/sqlx']} only %}

type {{ schema.Container.Name }} struct {
{% for repo in schema.Repositories %}	{{ repo.Interface.Name }} {{ spaces_after(repo.Interface.Name, schema.RepoInterfaceNameMaxLength) }}{{ repo.Interface.Call(_file.Package) }}{% if loop.last == false %}
{% endif %}{% endfor %}
}

func New{{ schema.Container.Name }}(db *sqlx.DB) *{{ schema.Container.Name }} {
	return &{{ schema.Container.Name }}{
{% for repo in schema.Repositories %}		{{ repo.Interface.Name }}: {{ spaces_after(repo.Name, schema.RepoNameMaxLength) }}New{{ repo.Name }}(db),{% if loop.last == false %}
{% endif %}{% endfor %}
	}
}
