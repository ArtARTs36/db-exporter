{% include '@embed/go-entities/repository_interfaces.go.tpl' with {'repositories': repositories, '_file': _file} only %}

{% for repo in repositories %}{% include '@embed/go-entities/repository_filters.go.tpl' with {'filters': repo.Filters, '_file': _file} only %}{% if loop.last == false %}
{% endif %}{% endfor %}