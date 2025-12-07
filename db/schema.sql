CREATE TABLE customers (
    id CHAR(36) PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(50) NOT NULL,
    balance_cents INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE rooms (
    room_number VARCHAR(10) PRIMARY KEY,
    type VARCHAR(50) NOT NULL
);

CREATE TABLE reservations (
    id CHAR(36) PRIMARY KEY,
    customer_id CHAR(36) NOT NULL,
    room_number VARCHAR(10) NOT NULL,
    check_in_date DATE NOT NULL,
    nights INT NOT NULL,
    status VARCHAR(20) NOT NULL,
    total_amount_cents INT NOT NULL,
    paid_amount_cents INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (room_number) REFERENCES rooms(room_number)
);

CREATE INDEX idx_reservations_room_date ON reservations(room_number, check_in_date);