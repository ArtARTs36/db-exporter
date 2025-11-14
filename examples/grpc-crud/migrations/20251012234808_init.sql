-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION "uuid-ossp";

CREATE TYPE user_status AS ENUM ('active', 'blocked');

CREATE TABLE users
(
    id           uuid NOT NULL PRIMARY KEY default uuid_generate_v4(),
    first_name   character varying NOT NULL,
    middle_name  character varying,
    created_at   timestamp without time zone NOT NULL default now(),
    status       user_status NOT NULL,
    deleted_at   timestamp without time zone
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TYPE user_status;
-- +goose StatementEnd
