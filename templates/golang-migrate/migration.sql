-- +migrate Up{% for query in up_queries %}
{{ query }}
{% endfor %}

-- +migrate Down{% for query in down_queries %}
{{ query }}
{% endfor %}
