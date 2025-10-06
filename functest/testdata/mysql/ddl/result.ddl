CREATE TABLE users
(
    `id`            int NOT NULL AUTO_INCREMENT,
    `username`      varchar(50) NOT NULL,
    `password_hash` varchar(255) NOT NULL,
    `created_at`    timestamp DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT users_id_pk PRIMARY KEY (`id`),
    CONSTRAINT users_username_uk UNIQUE (`username`)
);