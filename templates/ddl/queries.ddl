{% for query in migration.UpQueries %}{{ query }}{% if loop.last == false %}

{% endif %}{% endfor %}