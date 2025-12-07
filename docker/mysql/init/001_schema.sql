-- On utilise le nom de base 'xyzhotel' pour coller à ton main.go
CREATE DATABASE IF NOT EXISTS xyzhotel CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE xyzhotel;

-- 1. Table Customers
-- Changement: wallet_sold (DECIMAL) -> balance_cents (INT)
CREATE TABLE IF NOT EXISTS customers (
                                         id           CHAR(36)     NOT NULL PRIMARY KEY, -- UUID fixe
    email        VARCHAR(255) NOT NULL UNIQUE,
    phone        VARCHAR(50)  NOT NULL,
    full_name    VARCHAR(255) NOT NULL,
    balance_cents INT         NOT NULL DEFAULT 0,   -- Stocké en centimes
    created_at   TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 2. Table Rooms
-- Nécessaire pour les clés étrangères et la vérification d'existence
CREATE TABLE IF NOT EXISTS rooms (
                                     room_number VARCHAR(10) PRIMARY KEY, -- Ex: '101'
    type        VARCHAR(50) NOT NULL     -- 'STANDARD', 'SUPERIOR', 'SUITE'
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3. Table Reservations
-- Contient la logique de disponibilité et de paiement
CREATE TABLE IF NOT EXISTS reservations (
                                            id                 CHAR(36)    NOT NULL PRIMARY KEY,
    customer_id        CHAR(36)    NOT NULL,
    room_number        VARCHAR(10) NOT NULL,

    check_in_date      DATE        NOT NULL,
    nights             INT         NOT NULL,
    status             VARCHAR(20) NOT NULL, -- 'PENDING', 'CONFIRMED', 'CANCELLED'

    total_amount_cents INT         NOT NULL, -- Snapshot du prix total
    paid_amount_cents  INT         NOT NULL, -- Ce qui a été payé

    created_at         TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP   DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- Contraintes de clés étrangères
    CONSTRAINT fk_reservations_customer FOREIGN KEY (customer_id) REFERENCES customers(id),
    CONSTRAINT fk_reservations_room     FOREIGN KEY (room_number) REFERENCES rooms(room_number)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Index pour optimiser la vérification de disponibilité (IsRoomAvailable)
CREATE INDEX idx_reservations_room_date ON reservations(room_number, check_in_date);