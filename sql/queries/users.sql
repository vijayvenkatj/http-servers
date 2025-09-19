-- name: GetUsers :many
SELECT * FROM users ORDER BY created_at;

-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: DeleteUser :one
DELETE FROM users WHERE id = $1 RETURNING *;