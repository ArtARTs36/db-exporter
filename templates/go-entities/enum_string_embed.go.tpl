type {{ enum.Name.Value }} string

const (
{% for value in enum.Values %}    {{ value.Name }} {{ enum.Name.Value }} = "{{ value.Value }}"{% if loop.last == false %}
{% endif %}{% endfor %}
)