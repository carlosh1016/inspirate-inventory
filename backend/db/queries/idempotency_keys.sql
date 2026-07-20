-- name: GetIdempotencyKey :one
SELECT * FROM idempotency_keys
WHERE key = $1 AND expires_at > NOW();

-- name: InsertIdempotencyKey :exec
INSERT INTO idempotency_keys (key, usuario_id, endpoint, response_status, response_body, request_hash, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: DeleteExpiredIdempotencyKeys :exec
DELETE FROM idempotency_keys WHERE expires_at < NOW();
