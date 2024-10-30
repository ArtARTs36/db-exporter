type {{ enum.Name.Value }} string

const (
{% for value in enum.Values %}    {{ value.Name }} {{ enum.Name.Value }} = "{{ value.Value }}"{% if loop.last == false %}
{% endif %}{% endfor %}
)

func (e {{ enum.Name.Value }}) Valid() bool {
    switch e {
{% for value in enum.Values %}    case {{ value.Name }}:
        return true{% if loop.last == false %}
{% endif %}{% endfor %}
    default:
        return false
    }
}