{% for enum in enums %}{% include '@embed/go-entities/enum_string_embed.go.tpl' with {'enum': enum} only %}{% if loop.last == false %}

{% endif %}{% endfor %}