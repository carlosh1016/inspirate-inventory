-- name: InsertRefreshToken :one
INSERT INTO refresh_tokens (usuario_id, token_hash, ip_origen, user_agent, expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetRefreshTokenByHash :one
SELECT * FROM refresh_tokens WHERE token_hash = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1;

-- name: RevokeAllRefreshTokensByUser :exec
UPDATE refresh_tokens SET revoked_at = NOW()
WHERE usuario_id = $1 AND revoked_at IS NULL;
