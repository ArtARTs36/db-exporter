type {{ filters.List.Name }} struct {
{% for prop in filters.List.Properties.List %}	{{ prop.PluralName }} {{ spaces_after(prop.Name, entity.Properties.MaxPropPluralNameLength) }}[]{{ prop.Type }}{% if loop.last == false %}
{% endif %}{% endfor %}
}

type {{ filters.Get.Name }} struct {
{% for prop in filters.Get.Properties.List %}	{{ prop.Name }} {{ spaces_after(prop.Name, entity.Properties.MaxPropNameLength) }}{{ prop.Type }}{% if loop.last == false %}
{% endif %}{% endfor %}
}

type {{ filters.Delete.Name }} struct {
{% for prop in filters.Delete.Properties.List %}	{{ prop.PluralName }} {{ spaces_after(prop.Name, entity.Properties.MaxPropNameLength) }}[]{{ prop.Type }}{% if loop.last == false %}
{% endif %}{% endfor %}
}