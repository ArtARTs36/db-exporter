package {{ _file.Package.Name }}

{% include 'go-entities/enum_string_embed.go.tpl' with {'enum': enum, '_file': _file} only %}
