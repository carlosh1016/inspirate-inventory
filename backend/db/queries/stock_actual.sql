-- name: InsertStockActual :exec
INSERT INTO stock_actual (sede_id, tipo_item, item_id, ubicacion, cantidad)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (sede_id, tipo_item, item_id, ubicacion) DO NOTHING;

-- name: GetStockTotalByItem :one
SELECT
  COALESCE(SUM(CASE WHEN ubicacion = 'vitrina' THEN cantidad ELSE 0 END), 0)::numeric AS vitrina,
  COALESCE(SUM(CASE WHEN ubicacion = 'bodega' THEN cantidad ELSE 0 END), 0)::numeric AS bodega
FROM stock_actual
WHERE sede_id = $1 AND tipo_item = $2 AND item_id = $3;

-- name: GetStockActualForUpdate :one
-- Locks the row for the duration of the caller's transaction so concurrent
-- movimientos against the same (sede, tipo_item, item, ubicacion) serialize
-- instead of racing. Relies on InitializeStock having already created the
-- vitrina+bodega rows when the item itself was created.
SELECT * FROM stock_actual
WHERE sede_id = $1 AND tipo_item = $2 AND item_id = $3 AND ubicacion = $4
FOR UPDATE;

-- name: UpsertStockActual :one
INSERT INTO stock_actual (sede_id, tipo_item, item_id, ubicacion, cantidad)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (sede_id, tipo_item, item_id, ubicacion)
DO UPDATE SET cantidad = EXCLUDED.cantidad, updated_at = NOW()
RETURNING *;

-- name: ListStockUnificado :many
WITH stock_agregado AS (
  SELECT
    tipo_item,
    item_id,
    SUM(CASE WHEN ubicacion = 'vitrina' THEN cantidad ELSE 0 END) AS stock_vitrina,
    SUM(CASE WHEN ubicacion = 'bodega' THEN cantidad ELSE 0 END) AS stock_bodega,
    SUM(cantidad) AS stock_total
  FROM stock_actual sa_raw
  WHERE sa_raw.sede_id = @sede_id::bigint
  GROUP BY tipo_item, item_id
)
SELECT * FROM (
  -- Fragancias
  SELECT
    'fragancia'::text AS tipo_item,
    f.id AS item_id,
    f.nombre_comercial AS nombre,
    CASE WHEN f.nombre_alternativo IS NOT NULL
      THEN CONCAT(f.nombre_alternativo, ' (', f.genero::text, ')')
      ELSE f.genero::text
    END AS detalle_extra,
    COALESCE(sa.stock_vitrina, 0)::numeric AS stock_vitrina,
    COALESCE(sa.stock_bodega, 0)::numeric AS stock_bodega,
    COALESCE(sa.stock_total, 0)::numeric AS stock_total,
    f.gramos_minimo AS minimo,
    (COALESCE(sa.stock_total, 0) < f.gramos_minimo) AS bajo_minimo,
    'gramos'::text AS unidad
  FROM fragancias f
  LEFT JOIN stock_agregado sa ON sa.tipo_item = 'fragancia' AND sa.item_id = f.id
  WHERE f.deleted_at IS NULL
    AND (@include_inactivos::bool OR f.activo = true)
    AND (@tipo_item_filter::text = '' OR @tipo_item_filter::text = 'fragancia')

  UNION ALL

  -- Variantes de envase
  SELECT
    'variante_envase'::text AS tipo_item,
    ve.id AS item_id,
    CONCAT(me.tipo, ' ', me.tamano_oz::text, 'oz ', ve.color) AS nombre,
    CONCAT('Precio con fragancia: ', me.precio_con_fragancia::text) AS detalle_extra,
    COALESCE(sa.stock_vitrina, 0)::numeric AS stock_vitrina,
    COALESCE(sa.stock_bodega, 0)::numeric AS stock_bodega,
    COALESCE(sa.stock_total, 0)::numeric AS stock_total,
    ve.stock_minimo::numeric AS minimo,
    (COALESCE(sa.stock_total, 0) < ve.stock_minimo) AS bajo_minimo,
    'unidades'::text AS unidad
  FROM variantes_envase ve
  INNER JOIN modelos_envase me ON me.id = ve.modelo_envase_id
  LEFT JOIN stock_agregado sa ON sa.tipo_item = 'variante_envase' AND sa.item_id = ve.id
  WHERE ve.deleted_at IS NULL
    AND (@include_inactivos::bool OR ve.activo = true)
    AND (@tipo_item_filter::text = '' OR @tipo_item_filter::text = 'variante_envase')

  UNION ALL

  -- Productos
  SELECT
    'producto'::text AS tipo_item,
    p.id AS item_id,
    p.nombre,
    p.categoria::text AS detalle_extra,
    COALESCE(sa.stock_vitrina, 0)::numeric AS stock_vitrina,
    COALESCE(sa.stock_bodega, 0)::numeric AS stock_bodega,
    COALESCE(sa.stock_total, 0)::numeric AS stock_total,
    p.stock_minimo::numeric AS minimo,
    (COALESCE(sa.stock_total, 0) < p.stock_minimo) AS bajo_minimo,
    'unidades'::text AS unidad
  FROM productos p
  LEFT JOIN stock_agregado sa ON sa.tipo_item = 'producto' AND sa.item_id = p.id
  WHERE p.deleted_at IS NULL
    AND (@include_inactivos::bool OR p.activo = true)
    AND (@tipo_item_filter::text = '' OR @tipo_item_filter::text = 'producto')
) unified
WHERE
  (NOT @stock_bajo::bool OR bajo_minimo = true)
  AND (NOT @stock_cero::bool OR stock_total = 0)
ORDER BY tipo_item ASC, nombre ASC
LIMIT $1 OFFSET $2;

-- name: CountStockUnificado :one
-- Self-contained duplicate of ListStockUnificado's CTE/UNION/WHERE (same bug
-- class fixed in Tanda 2's CountFragancias: Count must mirror every filter
-- List applies, or meta.total drifts from the actually-returned items).
WITH stock_agregado AS (
  SELECT
    tipo_item,
    item_id,
    SUM(CASE WHEN ubicacion = 'vitrina' THEN cantidad ELSE 0 END) AS stock_vitrina,
    SUM(CASE WHEN ubicacion = 'bodega' THEN cantidad ELSE 0 END) AS stock_bodega,
    SUM(cantidad) AS stock_total
  FROM stock_actual sa_raw
  WHERE sa_raw.sede_id = @sede_id::bigint
  GROUP BY tipo_item, item_id
)
SELECT COUNT(*) FROM (
  SELECT
    f.id AS item_id,
    COALESCE(sa.stock_total, 0)::numeric AS stock_total,
    (COALESCE(sa.stock_total, 0) < f.gramos_minimo) AS bajo_minimo
  FROM fragancias f
  LEFT JOIN stock_agregado sa ON sa.tipo_item = 'fragancia' AND sa.item_id = f.id
  WHERE f.deleted_at IS NULL
    AND (@include_inactivos::bool OR f.activo = true)
    AND (@tipo_item_filter::text = '' OR @tipo_item_filter::text = 'fragancia')

  UNION ALL

  SELECT
    ve.id AS item_id,
    COALESCE(sa.stock_total, 0)::numeric AS stock_total,
    (COALESCE(sa.stock_total, 0) < ve.stock_minimo) AS bajo_minimo
  FROM variantes_envase ve
  INNER JOIN modelos_envase me ON me.id = ve.modelo_envase_id
  LEFT JOIN stock_agregado sa ON sa.tipo_item = 'variante_envase' AND sa.item_id = ve.id
  WHERE ve.deleted_at IS NULL
    AND (@include_inactivos::bool OR ve.activo = true)
    AND (@tipo_item_filter::text = '' OR @tipo_item_filter::text = 'variante_envase')

  UNION ALL

  SELECT
    p.id AS item_id,
    COALESCE(sa.stock_total, 0)::numeric AS stock_total,
    (COALESCE(sa.stock_total, 0) < p.stock_minimo) AS bajo_minimo
  FROM productos p
  LEFT JOIN stock_agregado sa ON sa.tipo_item = 'producto' AND sa.item_id = p.id
  WHERE p.deleted_at IS NULL
    AND (@include_inactivos::bool OR p.activo = true)
    AND (@tipo_item_filter::text = '' OR @tipo_item_filter::text = 'producto')
) unified
WHERE
  (NOT @stock_bajo::bool OR bajo_minimo = true)
  AND (NOT @stock_cero::bool OR stock_total = 0);
