-- name: CreateCustomer :exec
INSERT INTO customers (id, full_name, email, phone, balance_cents)
VALUES (?, ?, ?, ?, ?);

-- name: GetCustomerByID :one
SELECT * FROM customers
WHERE id = ? LIMIT 1;

-- name: GetCustomerByEmail :one
SELECT * FROM customers
WHERE email = ? LIMIT 1;

-- name: UpdateCustomer :exec
UPDATE customers
SET full_name = ?, email = ?, phone = ?, balance_cents = ?
WHERE id = ?;

-- name: ListCustomers :many
SELECT * FROM customers
ORDER BY created_at DESC;