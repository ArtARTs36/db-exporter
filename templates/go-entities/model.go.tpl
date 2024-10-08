package {{ package }}

{% if schema.Imports.Valid() %}import (
{% for im in schema.Imports.List() %}	"{{ im }}"{% if loop.last == false %}
{% endif %}{% endfor %}
){% endif %}
{% for table in schema.Tables %}
type {{ table.Name.Value }} struct {
{% for prop in table.Properties %}	{{ prop.Name }} {{ spaces(prop.NameOffset) }}{{ prop.Type }} {{ spaces(prop.TypeOffset) }}`db:"{{ prop.ColumnName }}"`{% if loop.last == false %}
{% endif %}{% endfor %}
}{% if loop.last == false %}
{% endif %}{% endfor %}
