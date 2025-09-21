-- name: StoreRefreshToken :exec

INSERT INTO refresh_tokens (token, user_id, created_at, updated_at, expires_at, revoked_at)
VALUES ($1, $2, NOW(), NOW(), $3, NULL)
ON CONFLICT (token) DO UPDATE
SET 
    user_id = EXCLUDED.user_id,
    updated_at = NOW(),
    expires_at = EXCLUDED.expires_at,
    revoked_at = NULL;

-- name: GetRefreshToken :one

SELECT * FROM refresh_tokens WHERE token = $1;

-- name: RevokeRefreshToken :exec

UPDATE refresh_tokens SET updated_at = NOW(), revoked_at = NOW() WHERE token = $1;

-- name: UpdateUser :one

UPDATE users SET email = $1, hashed_password = $2 WHERE id = $3 RETURNING *;