package models

{% if has_imports %}import (
{% for im in imports %}    "{{ im }}"{% if loop.last == false %}
{% endif %}{% endfor %}
){% endif %}
{% for table in schema.Tables %}
type {{ table.Name.Pascal() }} struct {
{% for prop in table.Properties %}	{{ prop.Name }} {{ spaces(prop.NameOffset) }}{{ prop.Type }} {{ spaces(prop.TypeOffset) }}`db:"{{ prop.ColumnName }}"` {% if loop.last == false %}
{% endif %}{% endfor %}
}{% if loop.last == false %}
{% endif %}{% endfor %}
