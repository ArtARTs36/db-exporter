type {{ entity.Name.Value }} struct {
{% for prop in entity.Properties %}	{{ prop.Name }} {{ spaces(prop.NameOffset) }}{{ prop.Type }} {{ spaces(prop.TypeOffset) }}`db:"{{ prop.ColumnName }}"`{% if loop.last == false %}
{% endif %}{% endfor %}
}