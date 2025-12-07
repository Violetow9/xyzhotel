CREATE DATABASE IF NOT EXISTS xyzhotel CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE xyzhotel;

CREATE TABLE IF NOT EXISTS customers (
    id           CHAR(36)     NOT NULL PRIMARY KEY,
    email        VARCHAR(255) NOT NULL UNIQUE,
    phone        VARCHAR(50)  NOT NULL,
    full_name    VARCHAR(255) NOT NULL,
    balance_cents INT         NOT NULL DEFAULT 0,
    created_at   TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS rooms (
    room_number VARCHAR(10) PRIMARY KEY,
    type        VARCHAR(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS reservations (
    id                 CHAR(36)    NOT NULL PRIMARY KEY,
    customer_id        CHAR(36)    NOT NULL,
    room_number        VARCHAR(10) NOT NULL,

    check_in_date      DATE        NOT NULL,
    nights             INT         NOT NULL,
    status             VARCHAR(20) NOT NULL,

    total_amount_cents INT         NOT NULL,
    paid_amount_cents  INT         NOT NULL,

    created_at         TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP   DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,


    CONSTRAINT fk_reservations_customer FOREIGN KEY (customer_id) REFERENCES customers(id),
    CONSTRAINT fk_reservations_room     FOREIGN KEY (room_number) REFERENCES rooms(room_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_reservations_room_date ON reservations(room_number, check_in_date);