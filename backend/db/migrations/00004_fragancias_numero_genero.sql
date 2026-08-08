-- Agrega numero_genero a fragancias: un código informativo por género
-- (masculina/femenina numeradas cada una desde 1, como las hojas de Excel
-- que ya usa la dueña del negocio). Es puramente una etiqueta visible — el
-- id sigue siendo la única FK real en venta_items, movimientos_inventario,
-- reportes y auditoría; nada de eso cambia.

-- +goose Up

-- +goose StatementBegin
ALTER TABLE fragancias ADD COLUMN numero_genero INT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
UPDATE fragancias f
SET numero_genero = sub.rn
FROM (
  SELECT id, ROW_NUMBER() OVER (PARTITION BY sede_id, genero ORDER BY id) AS rn
  FROM fragancias
) sub
WHERE f.id = sub.id;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE fragancias ALTER COLUMN numero_genero SET NOT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX uq_fragancias_sede_genero_numero ON fragancias (sede_id, genero, numero_genero) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP INDEX IF EXISTS uq_fragancias_sede_genero_numero;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE fragancias DROP COLUMN numero_genero;
-- +goose StatementEnd
