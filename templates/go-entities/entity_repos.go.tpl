{% include 'go-entities/go_file_header.go.tpl' with {'_file': _file, 'stdImports': ['context']} only %}

{% include 'go-entities/entity_repos_embed.go.tpl' with {'repositories': schema.Repositories, '_file': _file} only %}
