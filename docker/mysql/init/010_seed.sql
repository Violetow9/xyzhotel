USE xyzhotel;

INSERT INTO customers (id, email, phone, full_name, balance_cents)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'alice@example.com', '+33123456789', 'Alice Martin', 10000)
    ON DUPLICATE KEY UPDATE
        full_name = VALUES(full_name),
        balance_cents = VALUES(balance_cents);

INSERT INTO rooms (room_number, type) VALUES
                                          ('101', 'STANDARD'),
                                          ('102', 'STANDARD'),
                                          ('103', 'STANDARD'),
                                          ('201', 'SUPERIOR'),
                                          ('202', 'SUPERIOR'),
                                          ('301', 'SUITE')
    ON DUPLICATE KEY UPDATE type = VALUES(type);