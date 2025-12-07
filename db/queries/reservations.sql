-- name: CreateReservation :exec
INSERT INTO reservations (
    id, customer_id, room_number, check_in_date, nights, status,
    total_amount_cents, paid_amount_cents
) VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateReservation :exec
UPDATE reservations
SET status = ?, paid_amount_cents = ?
WHERE id = ?;

-- name: GetReservation :one
SELECT * FROM reservations
WHERE id = ? LIMIT 1;

-- name: ListReservations :many
SELECT * FROM reservations
ORDER BY check_in_date DESC;

-- name: IsRoomAvailable :one
SELECT COUNT(*)
FROM reservations
WHERE room_number = ?
  AND status != 'CANCELLED'
  AND (
    check_in_date < sqlc.arg(end_date)
    AND DATE_ADD(check_in_date, INTERVAL nights DAY) > sqlc.arg(start_date)
  );

-- name: GetReservationsByRoom :many
SELECT * FROM reservations
WHERE room_number = ?
ORDER BY check_in_date DESC;