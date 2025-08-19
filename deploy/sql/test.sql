CREATE TABLE `users` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `username` varchar(20) NOT NULL DEFAULT '',
    `email` varchar(100) NOT NULL DEFAULT '',
    `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
