-- name: GetCustomerByID :one
SELECT id, email, phone, full_name, wallet_sold
FROM customers
WHERE id = ? LIMIT 1;

-- name: GetCustomerByEmail :one
SELECT id, email, phone, full_name, wallet_sold
FROM customers
WHERE email = ? LIMIT 1;

-- name: ListCustomers :many
SELECT id, email, phone, full_name, wallet_sold
FROM customers
ORDER BY full_name ASC;

-- name: CreateCustomer :exec
INSERT INTO customers (id, email, phone, full_name, wallet_sold)
VALUES (?, ?, ?, ?, ?);

-- name: UpdateCustomer :execrows
UPDATE customers
SET email = ?, phone = ?, full_name = ?, wallet_sold = ?
WHERE id = ?;

-- name: DeleteCustomer :execrows
DELETE FROM customers
WHERE id = ?;