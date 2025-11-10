INSERT INTO customers (id, email, phone, full_name, wallet_sold)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'alice@example.com', '+33123456789', 'Alice Martin', 100.00)
    ON DUPLICATE KEY UPDATE full_name = VALUES(full_name);
