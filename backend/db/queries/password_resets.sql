-- name: InsertPasswordReset :one
INSERT INTO password_resets (usuario_id, token_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPasswordResetByHash :one
SELECT * FROM password_resets WHERE token_hash = $1;

-- name: MarkPasswordResetUsed :exec
UPDATE password_resets SET used_at = NOW() WHERE id = $1;
