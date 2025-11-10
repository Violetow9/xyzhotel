-- MySQL 8.4+
CREATE TABLE IF NOT EXISTS customers (
                                         id           VARCHAR(36)  NOT NULL PRIMARY KEY,        -- UUID string côté code
    email        VARCHAR(320) NOT NULL UNIQUE,              -- unique
    phone        VARCHAR(32)  NOT NULL,                     -- adapte si nullable
    full_name    VARCHAR(200) NOT NULL,
    wallet_sold  DECIMAL(18,2) NOT NULL DEFAULT 0.00       -- argent -> DECIMAL
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
