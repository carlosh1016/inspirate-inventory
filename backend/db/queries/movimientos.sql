-- name: ListMovimientosPaginated :many
SELECT m.*,
  (CASE m.tipo_item
    WHEN 'fragancia' THEN (SELECT nombre_comercial FROM fragancias WHERE id = m.item_id)
    WHEN 'variante_envase' THEN (
      SELECT CONCAT(me.tipo, ' ', me.tamano_oz::text, 'oz ', ve.color)
      FROM variantes_envase ve
      JOIN modelos_envase me ON me.id = ve.modelo_envase_id
      WHERE ve.id = m.item_id
    )
    WHEN 'producto' THEN (SELECT nombre FROM productos WHERE id = m.item_id)
  END)::text AS item_nombre,
  u.nombre_completo AS usuario_nombre
FROM movimientos_inventario m
INNER JOIN usuarios u ON u.id = m.usuario_id
WHERE
  (@sede_id::bigint = 0 OR m.sede_id = @sede_id)
  AND (@tipo_item::text = '' OR m.tipo_item::text = @tipo_item)
  AND (@item_id::bigint = 0 OR m.item_id = @item_id)
  AND (@tipo::text = '' OR m.tipo::text = @tipo)
  AND (@usuario_id::bigint = 0 OR m.usuario_id = @usuario_id)
  AND (@ubicacion::text = '' OR m.ubicacion::text = @ubicacion)
  AND (@venta_id::bigint = 0 OR m.venta_id = @venta_id)
  AND (@fecha_desde::timestamptz IS NULL OR m.created_at >= @fecha_desde)
  AND (@fecha_hasta::timestamptz IS NULL OR m.created_at <= @fecha_hasta)
ORDER BY m.created_at DESC, m.id DESC
LIMIT $1 OFFSET $2;

-- name: CountMovimientos :one
-- Mirrors ListMovimientosPaginated's filters exactly (no HAVING/aggregate
-- here, but the same principle applies: any filter added to List must be
-- added here too, or meta.total drifts from the actually-returned items).
SELECT COUNT(*) FROM movimientos_inventario m
WHERE
  (@sede_id::bigint = 0 OR m.sede_id = @sede_id)
  AND (@tipo_item::text = '' OR m.tipo_item::text = @tipo_item)
  AND (@item_id::bigint = 0 OR m.item_id = @item_id)
  AND (@tipo::text = '' OR m.tipo::text = @tipo)
  AND (@usuario_id::bigint = 0 OR m.usuario_id = @usuario_id)
  AND (@ubicacion::text = '' OR m.ubicacion::text = @ubicacion)
  AND (@venta_id::bigint = 0 OR m.venta_id = @venta_id)
  AND (@fecha_desde::timestamptz IS NULL OR m.created_at >= @fecha_desde)
  AND (@fecha_hasta::timestamptz IS NULL OR m.created_at <= @fecha_hasta);

-- name: InsertMovimiento :one
INSERT INTO movimientos_inventario (
  sede_id, usuario_id, tipo_item, item_id, tipo, ubicacion,
  cantidad, stock_anterior, stock_posterior, motivo, venta_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;
