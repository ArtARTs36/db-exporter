-- +goose Up
-- +goose StatementBegin
{% for query in migration.UpQueries %}
{{ query }}
{% endfor %}
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
{% for query in migration.DownQueries %}
{{ query }}
{% endfor %}
-- +goose StatementEnd
