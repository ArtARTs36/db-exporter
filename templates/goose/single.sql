-- +goose Up
-- +goose StatementBegin
{% for query in up_queries %}
{{ query }};
{% endfor %}
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
{% for query in down_queries %}
{{ query }};
{% endfor %}
-- +goose StatementEnd
