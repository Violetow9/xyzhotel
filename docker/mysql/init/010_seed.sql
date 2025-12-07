USE xyzhotel;

-- 1. Insertion du Client de test
-- Note: balance_cents = 10000 (pour 100.00€)
INSERT INTO customers (id, email, phone, full_name, balance_cents)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'alice@example.com', '+33123456789', 'Alice Martin', 10000)
    ON DUPLICATE KEY UPDATE
                         full_name = VALUES(full_name),
                         balance_cents = VALUES(balance_cents);

-- 2. Insertion des Chambres (Inventaire physique)
-- Sans cela, le Repository Room ne trouvera rien.
INSERT INTO rooms (room_number, type) VALUES
                                          ('101', 'STANDARD'),
                                          ('102', 'STANDARD'),
                                          ('103', 'STANDARD'),
                                          ('201', 'SUPERIOR'), -- Correspond à 'Deluxe' ou 'Superior' selon ton code
                                          ('202', 'SUPERIOR'),
                                          ('301', 'SUITE')
    ON DUPLICATE KEY UPDATE type = VALUES(type);