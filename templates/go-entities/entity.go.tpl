package {{ package.Name }}{% if schema.Imports.Valid() %}

import (
{% for im in schema.Imports.List() %}	"{{ im }}"{% if loop.last == false %}
{% endif %}{% endfor %}
){% endif %}{% if schema.Repositories | length > 0 %}
{% include 'go-entities/entity_repos.go.tpl' with {'repositories': schema.Repositories} only %}{% endif %}
{% for entity in schema.Entities %}
{% include 'go-entities/entity_struct.go.tpl' with {'entity': entity} only %}
{% endfor %}