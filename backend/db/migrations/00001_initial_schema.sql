-- Schema inicial completo de Inspírate Inventory: enums, tablas, índices y
-- triggers de updated_at. No inserta datos (la carga inicial es responsabilidad
-- de un módulo posterior de onboarding).

-- +goose Up

-- +goose StatementBegin
CREATE TYPE genero_enum AS ENUM ('masculina', 'femenina');
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TYPE categoria_producto_enum AS ENUM ('crema', 'feromona', 'hogar', 'regalo', 'otro');
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TYPE rol_enum AS ENUM ('admin', 'vendedora');
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TYPE tipo_item_enum AS ENUM ('fragancia', 'variante_envase', 'producto');
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TYPE ubicacion_enum AS ENUM ('vitrina', 'bodega');
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TYPE tipo_movimiento_enum AS ENUM (
  'entrada_mercancia',
  'traslado_bodega_vitrina',
  'venta',
  'ajuste',
  'danado',
  'devolucion',
  'correccion'
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TYPE tipo_linea_enum AS ENUM ('envase_con_fragancia', 'recarga', 'envase_solo', 'producto_otro');
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TYPE estado_cuadre_enum AS ENUM ('abierto', 'cerrado');
-- +goose StatementEnd

-- 1. sedes
-- +goose StatementBegin
CREATE TABLE sedes (
  id BIGSERIAL PRIMARY KEY,
  nombre TEXT NOT NULL,
  activa BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- 2. usuarios
-- +goose StatementBegin
CREATE TABLE usuarios (
  id BIGSERIAL PRIMARY KEY,
  sede_id BIGINT NOT NULL REFERENCES sedes(id) ON DELETE RESTRICT,
  nombre_completo TEXT NOT NULL,
  correo TEXT NOT NULL,
  password_hash TEXT NOT NULL,
  rol rol_enum NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT true,
  last_login_at TIMESTAMPTZ NULL,
  deleted_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX uq_usuarios_correo ON usuarios (LOWER(correo)) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- 3. refresh_tokens
-- +goose StatementBegin
CREATE TABLE refresh_tokens (
  id BIGSERIAL PRIMARY KEY,
  usuario_id BIGINT NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
  token_hash TEXT NOT NULL UNIQUE,
  ip_origen INET NULL,
  user_agent TEXT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- 4. password_resets
-- +goose StatementBegin
CREATE TABLE password_resets (
  id BIGSERIAL PRIMARY KEY,
  usuario_id BIGINT NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
  token_hash TEXT NOT NULL UNIQUE,
  expires_at TIMESTAMPTZ NOT NULL,
  used_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- 5. fragancias
-- +goose StatementBegin
CREATE TABLE fragancias (
  id BIGSERIAL PRIMARY KEY,
  sede_id BIGINT NOT NULL REFERENCES sedes(id) ON DELETE RESTRICT,
  nombre_comercial TEXT NOT NULL,
  nombre_alternativo TEXT NULL,
  genero genero_enum NOT NULL,
  gramos_minimo NUMERIC(10,2) NOT NULL DEFAULT 0 CHECK (gramos_minimo >= 0),
  activo BOOLEAN NOT NULL DEFAULT true,
  deleted_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- Único por sede: dos sedes distintas pueden vender una fragancia con el
-- mismo nombre comercial, pero no puede repetirse dentro de la misma sede.
-- +goose StatementBegin
CREATE UNIQUE INDEX uq_fragancias_sede_nombre ON fragancias (sede_id, LOWER(nombre_comercial)) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- 6. modelos_envase
-- +goose StatementBegin
CREATE TABLE modelos_envase (
  id BIGSERIAL PRIMARY KEY,
  tipo TEXT NOT NULL,
  tamano_oz NUMERIC(4,2) NOT NULL CHECK (tamano_oz > 0),
  equiv_gramos NUMERIC(10,2) NOT NULL CHECK (equiv_gramos > 0),
  precio_solo NUMERIC(12,2) NOT NULL CHECK (precio_solo >= 0),
  precio_con_fragancia NUMERIC(12,2) NOT NULL CHECK (precio_con_fragancia >= 0),
  precio_recarga NUMERIC(12,2) NOT NULL CHECK (precio_recarga >= 0),
  activo BOOLEAN NOT NULL DEFAULT true,
  deleted_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX uq_modelos_envase_tipo_tamano ON modelos_envase (LOWER(tipo), tamano_oz) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- 7. variantes_envase
-- +goose StatementBegin
CREATE TABLE variantes_envase (
  id BIGSERIAL PRIMARY KEY,
  modelo_envase_id BIGINT NOT NULL REFERENCES modelos_envase(id) ON DELETE RESTRICT,
  sede_id BIGINT NOT NULL REFERENCES sedes(id) ON DELETE RESTRICT,
  color TEXT NOT NULL,
  stock_minimo INTEGER NOT NULL DEFAULT 0 CHECK (stock_minimo >= 0),
  activo BOOLEAN NOT NULL DEFAULT true,
  deleted_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX uq_variantes_envase_modelo_color ON variantes_envase (modelo_envase_id, LOWER(color)) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- 8. productos
-- +goose StatementBegin
CREATE TABLE productos (
  id BIGSERIAL PRIMARY KEY,
  sede_id BIGINT NOT NULL REFERENCES sedes(id) ON DELETE RESTRICT,
  nombre TEXT NOT NULL,
  categoria categoria_producto_enum NOT NULL,
  precio NUMERIC(12,2) NOT NULL CHECK (precio >= 0),
  stock_minimo INTEGER NOT NULL DEFAULT 0 CHECK (stock_minimo >= 0),
  activo BOOLEAN NOT NULL DEFAULT true,
  deleted_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX uq_productos_sede_nombre_categoria ON productos (sede_id, LOWER(nombre), categoria) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- 9. stock_actual: snapshot polimórfico (tipo_item + item_id). No lleva FK
-- real hacia fragancias/variantes_envase/productos porque un mismo item_id
-- puede corresponder a tablas distintas según tipo_item; la integridad
-- referencial de este campo se garantiza a nivel de aplicación.
-- +goose StatementBegin
CREATE TABLE stock_actual (
  id BIGSERIAL PRIMARY KEY,
  sede_id BIGINT NOT NULL REFERENCES sedes(id) ON DELETE RESTRICT,
  tipo_item tipo_item_enum NOT NULL,
  item_id BIGINT NOT NULL,
  ubicacion ubicacion_enum NOT NULL,
  cantidad NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (cantidad >= 0),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (sede_id, tipo_item, item_id, ubicacion)
);
-- +goose StatementEnd

-- 10. movimientos_inventario: fuente de verdad, filas inmutables.
-- venta_id no lleva FK todavía porque ventas se crea después; se agrega con
-- un ALTER diferido más abajo.
-- +goose StatementBegin
CREATE TABLE movimientos_inventario (
  id BIGSERIAL PRIMARY KEY,
  sede_id BIGINT NOT NULL REFERENCES sedes(id) ON DELETE RESTRICT,
  usuario_id BIGINT NOT NULL REFERENCES usuarios(id) ON DELETE RESTRICT,
  tipo_item tipo_item_enum NOT NULL,
  item_id BIGINT NOT NULL,
  tipo tipo_movimiento_enum NOT NULL,
  ubicacion ubicacion_enum NOT NULL,
  cantidad NUMERIC(12,2) NOT NULL,
  stock_anterior NUMERIC(12,2) NOT NULL CHECK (stock_anterior >= 0),
  stock_posterior NUMERIC(12,2) NOT NULL CHECK (stock_posterior >= 0),
  motivo TEXT NULL,
  venta_id BIGINT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT chk_movimientos_motivo CHECK (
    tipo NOT IN ('ajuste', 'danado', 'correccion')
    OR (motivo IS NOT NULL AND length(trim(motivo)) > 0)
  )
);
-- +goose StatementEnd

-- 11. metodos_pago
-- +goose StatementBegin
CREATE TABLE metodos_pago (
  id BIGSERIAL PRIMARY KEY,
  nombre TEXT NOT NULL,
  codigo TEXT NOT NULL,
  activo BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX uq_metodos_pago_codigo ON metodos_pago (LOWER(codigo));
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX uq_metodos_pago_nombre ON metodos_pago (LOWER(nombre));
-- +goose StatementEnd

-- 12. ventas
-- +goose StatementBegin
CREATE TABLE ventas (
  id BIGSERIAL PRIMARY KEY,
  sede_id BIGINT NOT NULL REFERENCES sedes(id) ON DELETE RESTRICT,
  usuario_id BIGINT NOT NULL REFERENCES usuarios(id) ON DELETE RESTRICT,
  metodo_pago_id BIGINT NOT NULL REFERENCES metodos_pago(id) ON DELETE RESTRICT,
  subtotal NUMERIC(12,2) NOT NULL CHECK (subtotal >= 0),
  descuento_pct NUMERIC(5,2) NOT NULL DEFAULT 0 CHECK (descuento_pct >= 0 AND descuento_pct <= 100),
  descuento_monto NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (descuento_monto >= 0),
  total NUMERIC(12,2) NOT NULL CHECK (total >= 0),
  observaciones TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- FK diferida: movimientos_inventario.venta_id -> ventas.id. Ventas nunca se
-- borran, así que RESTRICT (no CASCADE) es la semántica correcta.
-- +goose StatementBegin
ALTER TABLE movimientos_inventario
  ADD CONSTRAINT fk_movimientos_venta FOREIGN KEY (venta_id) REFERENCES ventas(id) ON DELETE RESTRICT;
-- +goose StatementEnd

-- 13. venta_items: líneas polimórficas; los FK requeridos varían según tipo_linea.
-- +goose StatementBegin
CREATE TABLE venta_items (
  id BIGSERIAL PRIMARY KEY,
  venta_id BIGINT NOT NULL REFERENCES ventas(id) ON DELETE RESTRICT,
  tipo_linea tipo_linea_enum NOT NULL,
  fragancia_id BIGINT NULL REFERENCES fragancias(id) ON DELETE RESTRICT,
  variante_envase_id BIGINT NULL REFERENCES variantes_envase(id) ON DELETE RESTRICT,
  producto_id BIGINT NULL REFERENCES productos(id) ON DELETE RESTRICT,
  feromona_producto_id BIGINT NULL REFERENCES productos(id) ON DELETE RESTRICT,
  gramos_fragancia NUMERIC(10,2) NULL CHECK (gramos_fragancia IS NULL OR gramos_fragancia > 0),
  cantidad INTEGER NOT NULL CHECK (cantidad > 0),
  precio_unitario NUMERIC(12,2) NOT NULL CHECK (precio_unitario >= 0),
  subtotal NUMERIC(12,2) NOT NULL CHECK (subtotal >= 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT chk_venta_items_tipo_linea CHECK (
    (tipo_linea = 'envase_con_fragancia' AND fragancia_id IS NOT NULL AND variante_envase_id IS NOT NULL AND producto_id IS NULL AND gramos_fragancia IS NOT NULL)
    OR
    (tipo_linea = 'recarga' AND fragancia_id IS NOT NULL AND variante_envase_id IS NOT NULL AND producto_id IS NULL AND gramos_fragancia IS NOT NULL)
    OR
    (tipo_linea = 'envase_solo' AND fragancia_id IS NULL AND variante_envase_id IS NOT NULL AND producto_id IS NULL AND gramos_fragancia IS NULL AND feromona_producto_id IS NULL)
    OR
    (tipo_linea = 'producto_otro' AND fragancia_id IS NULL AND variante_envase_id IS NULL AND producto_id IS NOT NULL AND gramos_fragancia IS NULL AND feromona_producto_id IS NULL)
  )
);
-- +goose StatementEnd

-- 14. cuadres_caja
-- +goose StatementBegin
CREATE TABLE cuadres_caja (
  id BIGSERIAL PRIMARY KEY,
  sede_id BIGINT NOT NULL REFERENCES sedes(id) ON DELETE RESTRICT,
  fecha DATE NOT NULL,
  estado estado_cuadre_enum NOT NULL DEFAULT 'abierto',
  fondo_base NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (fondo_base >= 0),
  total_efectivo NUMERIC(12,2) NOT NULL DEFAULT 0,
  total_nequi NUMERIC(12,2) NOT NULL DEFAULT 0,
  total_daviplata NUMERIC(12,2) NOT NULL DEFAULT 0,
  total_transferencia NUMERIC(12,2) NOT NULL DEFAULT 0,
  total_otros NUMERIC(12,2) NOT NULL DEFAULT 0,
  total_pagos NUMERIC(12,2) NOT NULL DEFAULT 0,
  total_consignaciones NUMERIC(12,2) NOT NULL DEFAULT 0,
  valor_turno NUMERIC(12,2) NOT NULL DEFAULT 0,
  saldo_calculado NUMERIC(12,2) NOT NULL DEFAULT 0,
  observaciones TEXT NULL,
  cerrado_por_usuario_id BIGINT NULL REFERENCES usuarios(id) ON DELETE RESTRICT,
  cerrado_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (sede_id, fecha)
);
-- +goose StatementEnd

-- 15. pagos_caja
-- +goose StatementBegin
CREATE TABLE pagos_caja (
  id BIGSERIAL PRIMARY KEY,
  cuadre_caja_id BIGINT NOT NULL REFERENCES cuadres_caja(id) ON DELETE CASCADE,
  usuario_id BIGINT NOT NULL REFERENCES usuarios(id) ON DELETE RESTRICT,
  concepto TEXT NOT NULL,
  monto NUMERIC(12,2) NOT NULL CHECK (monto > 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- 16. consignaciones
-- +goose StatementBegin
CREATE TABLE consignaciones (
  id BIGSERIAL PRIMARY KEY,
  cuadre_caja_id BIGINT NOT NULL REFERENCES cuadres_caja(id) ON DELETE CASCADE,
  usuario_id BIGINT NOT NULL REFERENCES usuarios(id) ON DELETE RESTRICT,
  monto NUMERIC(12,2) NOT NULL CHECK (monto > 0),
  banco TEXT NULL,
  referencia TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- 17. sesiones_laborales
-- +goose StatementBegin
CREATE TABLE sesiones_laborales (
  id BIGSERIAL PRIMARY KEY,
  sede_id BIGINT NOT NULL REFERENCES sedes(id) ON DELETE RESTRICT,
  usuario_id BIGINT NOT NULL REFERENCES usuarios(id) ON DELETE RESTRICT,
  entrada_at TIMESTAMPTZ NOT NULL,
  salida_at TIMESTAMPTZ NULL,
  horas_trabajadas INTERVAL NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- 18. auditoria
-- +goose StatementBegin
CREATE TABLE auditoria (
  id BIGSERIAL PRIMARY KEY,
  usuario_id BIGINT NULL REFERENCES usuarios(id) ON DELETE RESTRICT,
  accion TEXT NOT NULL,
  tabla_afectada TEXT NULL,
  registro_id BIGINT NULL,
  datos_antes JSONB NULL,
  datos_despues JSONB NULL,
  ip INET NULL,
  user_agent TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- 19. idempotency_keys
-- +goose StatementBegin
CREATE TABLE idempotency_keys (
  id BIGSERIAL PRIMARY KEY,
  key TEXT NOT NULL UNIQUE,
  usuario_id BIGINT NOT NULL REFERENCES usuarios(id) ON DELETE RESTRICT,
  endpoint TEXT NOT NULL,
  response_status INTEGER NOT NULL,
  response_body JSONB NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- Índices explícitos

-- +goose StatementBegin
CREATE INDEX idx_usuarios_activos ON usuarios(is_active) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_fragancias_sede_activo ON fragancias(sede_id, activo) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_fragancias_genero ON fragancias(genero) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_variantes_modelo ON variantes_envase(modelo_envase_id) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_productos_categoria ON productos(categoria) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_stock_actual_lookup ON stock_actual(sede_id, tipo_item, item_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_movimientos_item ON movimientos_inventario(tipo_item, item_id, created_at DESC);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_movimientos_sede_fecha ON movimientos_inventario(sede_id, created_at DESC);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_movimientos_venta ON movimientos_inventario(venta_id) WHERE venta_id IS NOT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_movimientos_usuario ON movimientos_inventario(usuario_id, created_at DESC);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_ventas_sede_fecha ON ventas(sede_id, created_at DESC);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_ventas_usuario ON ventas(usuario_id, created_at DESC);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_ventas_fecha ON ventas(created_at);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_venta_items_venta ON venta_items(venta_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_cuadres_sede_estado ON cuadres_caja(sede_id, estado);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_pagos_cuadre ON pagos_caja(cuadre_caja_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_consignaciones_cuadre ON consignaciones(cuadre_caja_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_sesiones_usuario ON sesiones_laborales(usuario_id, entrada_at DESC);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_sesiones_abiertas ON sesiones_laborales(usuario_id) WHERE salida_at IS NULL;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_auditoria_usuario ON auditoria(usuario_id, created_at DESC);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_auditoria_accion ON auditoria(accion, created_at DESC);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_refresh_tokens_expires ON refresh_tokens(expires_at) WHERE revoked_at IS NULL;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_password_resets_expires ON password_resets(expires_at) WHERE used_at IS NULL;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_idempotency_expires ON idempotency_keys(expires_at);
-- +goose StatementEnd

-- Función y triggers de updated_at

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_sedes_updated_at BEFORE UPDATE ON sedes FOR EACH ROW EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_usuarios_updated_at BEFORE UPDATE ON usuarios FOR EACH ROW EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_fragancias_updated_at BEFORE UPDATE ON fragancias FOR EACH ROW EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_modelos_envase_updated_at BEFORE UPDATE ON modelos_envase FOR EACH ROW EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_variantes_envase_updated_at BEFORE UPDATE ON variantes_envase FOR EACH ROW EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_productos_updated_at BEFORE UPDATE ON productos FOR EACH ROW EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_metodos_pago_updated_at BEFORE UPDATE ON metodos_pago FOR EACH ROW EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_cuadres_caja_updated_at BEFORE UPDATE ON cuadres_caja FOR EACH ROW EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_sesiones_laborales_updated_at BEFORE UPDATE ON sesiones_laborales FOR EACH ROW EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose Down

-- Triggers
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_sesiones_laborales_updated_at ON sesiones_laborales;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_cuadres_caja_updated_at ON cuadres_caja;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_metodos_pago_updated_at ON metodos_pago;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_productos_updated_at ON productos;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_variantes_envase_updated_at ON variantes_envase;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_modelos_envase_updated_at ON modelos_envase;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_fragancias_updated_at ON fragancias;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_usuarios_updated_at ON usuarios;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_sedes_updated_at ON sedes;
-- +goose StatementEnd

-- Función
-- +goose StatementBegin
DROP FUNCTION IF EXISTS set_updated_at();
-- +goose StatementEnd

-- Índices explícitos
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_idempotency_expires;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_password_resets_expires;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_refresh_tokens_expires;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_auditoria_accion;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_auditoria_usuario;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_sesiones_abiertas;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_sesiones_usuario;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_consignaciones_cuadre;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_pagos_cuadre;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_cuadres_sede_estado;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_venta_items_venta;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_ventas_fecha;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_ventas_usuario;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_ventas_sede_fecha;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_movimientos_usuario;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_movimientos_venta;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_movimientos_sede_fecha;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_movimientos_item;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_stock_actual_lookup;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_productos_categoria;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_variantes_modelo;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_fragancias_genero;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_fragancias_sede_activo;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_usuarios_activos;
-- +goose StatementEnd

-- FK diferida
-- +goose StatementBegin
ALTER TABLE movimientos_inventario DROP CONSTRAINT IF EXISTS fk_movimientos_venta;
-- +goose StatementEnd

-- Tablas, en orden inverso de creación
-- +goose StatementBegin
DROP TABLE IF EXISTS idempotency_keys CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS auditoria CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS sesiones_laborales CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS consignaciones CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS pagos_caja CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS cuadres_caja CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS venta_items CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS ventas CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS metodos_pago CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS movimientos_inventario CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS stock_actual CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS productos CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS variantes_envase CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS modelos_envase CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS fragancias CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS password_resets CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS refresh_tokens CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS usuarios CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS sedes CASCADE;
-- +goose StatementEnd

-- Enums
-- +goose StatementBegin
DROP TYPE IF EXISTS estado_cuadre_enum;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TYPE IF EXISTS tipo_linea_enum;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TYPE IF EXISTS tipo_movimiento_enum;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TYPE IF EXISTS ubicacion_enum;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TYPE IF EXISTS tipo_item_enum;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TYPE IF EXISTS rol_enum;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TYPE IF EXISTS categoria_producto_enum;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TYPE IF EXISTS genero_enum;
-- +goose StatementEnd
