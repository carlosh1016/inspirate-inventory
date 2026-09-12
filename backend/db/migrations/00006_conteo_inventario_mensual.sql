-- +goose Up

-- +goose StatementBegin
CREATE TYPE estado_conteo_enum AS ENUM ('abierto', 'cerrado');
-- +goose StatementEnd

-- conteos_inventario: un conteo mensual por sede/mes. periodo siempre es el
-- primer día del mes (el CHECK evita crear conteos "a mitad de mes").
-- +goose StatementBegin
CREATE TABLE conteos_inventario (
  id BIGSERIAL PRIMARY KEY,
  sede_id BIGINT NOT NULL REFERENCES sedes(id) ON DELETE RESTRICT,
  periodo DATE NOT NULL CHECK (periodo = date_trunc('month', periodo)::date),
  estado estado_conteo_enum NOT NULL DEFAULT 'abierto',
  creado_por_usuario_id BIGINT NOT NULL REFERENCES usuarios(id) ON DELETE RESTRICT,
  cerrado_por_usuario_id BIGINT NULL REFERENCES usuarios(id) ON DELETE RESTRICT,
  cerrado_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (sede_id, periodo)
);
-- +goose StatementEnd

-- conteo_inventario_items: una fila por fragancia dentro de un conteo.
-- gramos_fisico queda NULL hasta que se cuenta físicamente esa fragancia.
-- +goose StatementBegin
CREATE TABLE conteo_inventario_items (
  id BIGSERIAL PRIMARY KEY,
  conteo_id BIGINT NOT NULL REFERENCES conteos_inventario(id) ON DELETE CASCADE,
  fragancia_id BIGINT NOT NULL REFERENCES fragancias(id) ON DELETE RESTRICT,
  saldo_inicial NUMERIC(10,2) NOT NULL CHECK (saldo_inicial >= 0),
  gramos_sistema NUMERIC(10,2) NOT NULL CHECK (gramos_sistema >= 0),
  gramos_fisico NUMERIC(10,2) NULL CHECK (gramos_fisico IS NULL OR gramos_fisico >= 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (conteo_id, fragancia_id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_conteos_inventario_sede ON conteos_inventario (sede_id, periodo DESC);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_conteo_items_conteo ON conteo_inventario_items (conteo_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_conteos_inventario_updated_at BEFORE UPDATE ON conteos_inventario FOR EACH ROW EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_conteo_items_updated_at BEFORE UPDATE ON conteo_inventario_items FOR EACH ROW EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_conteo_items_updated_at ON conteo_inventario_items;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_conteos_inventario_updated_at ON conteos_inventario;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS conteo_inventario_items CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS conteos_inventario CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TYPE IF EXISTS estado_conteo_enum;
-- +goose StatementEnd
