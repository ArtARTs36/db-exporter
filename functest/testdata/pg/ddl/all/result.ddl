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

    CONSTRAINT users_pk PRIMARY KEY (id)
);

CREATE TABLE phones
(
    user_id integer NOT NULL,
    number  character varying NOT NULL,

    CONSTRAINT phones_pk PRIMARY KEY (user_id, number),
    CONSTRAINT phone_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id)
);