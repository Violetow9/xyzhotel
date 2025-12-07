-- name: GetRoom :one
SELECT * FROM rooms
WHERE room_number = ? LIMIT 1;

-- name: ListRooms :many
SELECT * FROM rooms
ORDER BY room_number;

-- name: CreateRoom :exec
INSERT INTO rooms (room_number, type)
VALUES (?, ?);