{% for enum in enums %}{% include 'go-entities/enum_string_embed.go.tpl' with {'enum': enum} only %}{% if loop.last == false %}

{% endif %}{% endfor %}