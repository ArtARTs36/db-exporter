CREATE TABLE phones
(
    user_id integer NOT NULL,
    number  character varying NOT NULL,

    CONSTRAINT phones_pk PRIMARY KEY (user_id, number),
    CONSTRAINT phone_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id)
);