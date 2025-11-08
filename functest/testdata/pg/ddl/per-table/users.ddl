CREATE TYPE mood AS ENUM ('ok', 'happy');

CREATE TABLE users
(
    id           integer NOT NULL,
    name         character varying(64) NOT NULL,
    balance      real NOT NULL,
    prev_balance real,
    created_at   timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    current_mood mood NOT NULL,
    updated_at   timestamp without time zone,
    deleted_at   timestamp without time zone,

    CONSTRAINT users_pk PRIMARY KEY (id)
);

COMMENT ON COLUMN users.name IS 'user name';