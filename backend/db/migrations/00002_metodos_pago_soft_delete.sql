-- Agrega soft delete a metodos_pago: si un método de pago tiene ventas
-- asociadas no se puede borrar físicamente (violaría la FK de ventas), así
-- que a partir de esta migración se puede desactivar/ocultar en su lugar.

-- +goose Up

-- +goose StatementBegin
ALTER TABLE metodos_pago ADD COLUMN deleted_at TIMESTAMPTZ NULL;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS uq_metodos_pago_codigo;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS uq_metodos_pago_nombre;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX uq_metodos_pago_codigo ON metodos_pago (LOWER(codigo)) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX uq_metodos_pago_nombre ON metodos_pago (LOWER(nombre)) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP INDEX IF EXISTS uq_metodos_pago_codigo;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS uq_metodos_pago_nombre;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX uq_metodos_pago_codigo ON metodos_pago (LOWER(codigo));
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX uq_metodos_pago_nombre ON metodos_pago (LOWER(nombre));
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE metodos_pago DROP COLUMN deleted_at;
-- +goose StatementEnd
