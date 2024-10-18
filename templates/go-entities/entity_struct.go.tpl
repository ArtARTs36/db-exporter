type {{ entity.Name.Value }} struct {
{% for prop in entity.Properties.List %}	{{ prop.Name.Value }} {{ spaces_after(prop.Name, entity.Properties.MaxPropNameLength) }}{{ prop.Type }} {{ spaces_after(prop.Type, entity.Properties.MaxTypeNameLength) }}`db:"{{ prop.Column.Name.Value }}"`{% if loop.last == false %}
{% endif %}{% endfor %}
}