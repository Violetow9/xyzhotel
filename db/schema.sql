CREATE TABLE customers (
    id CHAR(36) PRIMARY KEY,       -- UUID stocké en string
                           full_name VARCHAR(255) NOT NULL,
                           email VARCHAR(255) NOT NULL UNIQUE,
                           phone VARCHAR(50) NOT NULL,
                           balance_cents INT NOT NULL DEFAULT 0, -- Le Wallet
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE rooms (
                       room_number VARCHAR(10) PRIMARY KEY, -- Ex: '101', '204B'
                       type VARCHAR(50) NOT NULL            -- 'STANDARD', 'SUPERIOR', 'SUITE'
    -- Note: Le prix est géré dans le code Go (Config), mais on pourrait le mettre ici aussi.
);

CREATE TABLE reservations (
                              id CHAR(36) PRIMARY KEY,
                              customer_id CHAR(36) NOT NULL,
                              room_number VARCHAR(10) NOT NULL,

                              check_in_date DATE NOT NULL,
                              nights INT NOT NULL,

                              status VARCHAR(20) NOT NULL, -- 'PENDING', 'CONFIRMED', 'CANCELLED'

                              total_amount_cents INT NOT NULL,
                              paid_amount_cents INT NOT NULL,

                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

                              FOREIGN KEY (customer_id) REFERENCES customers(id),
                              FOREIGN KEY (room_number) REFERENCES rooms(room_number)
);

-- Index pour la performance des recherches de disponibilité
CREATE INDEX idx_reservations_room_date ON reservations(room_number, check_in_date);