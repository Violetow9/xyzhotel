CREATE DATABASE IF NOT EXISTS app CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE app;

CREATE TABLE IF NOT EXISTS customers (
    id           VARCHAR(36)  NOT NULL PRIMARY KEY,
    email        VARCHAR(320) NOT NULL UNIQUE,
    phone        VARCHAR(32)  NOT NULL,
    full_name    VARCHAR(200) NOT NULL,
    wallet_sold  DECIMAL(18,2) NOT NULL DEFAULT 0.00
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
