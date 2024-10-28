-- +migrate Up{% for query in migration.UpQueries %}
{{ query }}
{% endfor %}

-- +migrate Down{% for query in migration.DownQueries %}
{{ query }}
{% endfor %}
