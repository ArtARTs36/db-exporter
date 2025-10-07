CREATE TABLE users
(
    `id`            int NOT NULL AUTO_INCREMENT,
    `username`      varchar(50) NOT NULL,
    `password_hash` varchar(255) NOT NULL,
    `created_at`    timestamp DEFAULT CURRENT_TIMESTAMP,
    `status`        ENUM('active', 'banned'),

    CONSTRAINT `PRIMARY` PRIMARY KEY (`id`),
    CONSTRAINT `username` UNIQUE (`username`)
);

CREATE TABLE orders
(
    `id`      int NOT NULL AUTO_INCREMENT,
    `user_id` int,

    CONSTRAINT `PRIMARY` PRIMARY KEY (`id`),
    CONSTRAINT orders_user_id_fk FOREIGN KEY (`user_id`) REFERENCES users (`id`)
);