{% include 'go-entities/repository_interfaces.go.tpl' with {'repositories': repositories} only %}
{% for repo in repositories %}
{% include 'go-entities/repository_filters.go.tpl' with {'filters': repo.Filters} only %}{% endfor %}