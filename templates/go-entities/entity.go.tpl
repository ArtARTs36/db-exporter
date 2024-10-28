{% set stdImports = [] %}{% if schema.Repositories | length > 0 %}{% if schema.WithMocks %}{% include 'go-entities/gomock_call.go.tpl' with {'_file': _file, 'repositories': schema.Repositories} only %}{% endif %}
{% set stdImports = ['context'] %}{% endif %}{% include 'go-entities/go_file_header.go.tpl' with {'_file': _file, 'stdImports': stdImports} only %}{% if schema.Repositories | length > 0 %}

{% include 'go-entities/entity_repos_embed.go.tpl' with {'repositories': schema.Repositories, '_file': _file} only %}{% endif %}
{% for entity in schema.Entities %}
{% include 'go-entities/entity_struct.go.tpl' with {'entity': entity} only %}
{% endfor %}