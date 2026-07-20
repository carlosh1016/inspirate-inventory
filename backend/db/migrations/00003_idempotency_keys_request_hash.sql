-- Agrega request_hash a idempotency_keys: la tabla original (M1) solo
-- guarda la respuesta cacheada, no un fingerprint del request original, así
-- que no hay forma de detectar que la misma key se reusó con un body
-- distinto. Con este hash (SHA-256 hex del body JSON crudo) POST /ventas
-- puede distinguir "repetición legítima" (mismo body) de "reuso indebido de
-- la key" (body distinto) y responder 409 en el segundo caso.

-- +goose Up

-- +goose StatementBegin
ALTER TABLE idempotency_keys ADD COLUMN request_hash TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE idempotency_keys ALTER COLUMN request_hash DROP DEFAULT;
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
ALTER TABLE idempotency_keys DROP COLUMN request_hash;
-- +goose StatementEnd
