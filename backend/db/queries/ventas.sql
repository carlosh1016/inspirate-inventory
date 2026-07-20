-- name: InsertVenta :one
INSERT INTO ventas (
    sede_id, usuario_id, metodo_pago_id,
    subtotal, descuento_pct, descuento_monto, total,
    observaciones
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetVentaByID :one
SELECT v.*, mp.nombre AS metodo_pago_nombre, mp.codigo AS metodo_pago_codigo,
       u.nombre_completo AS usuario_nombre
FROM ventas v
INNER JOIN metodos_pago mp ON mp.id = v.metodo_pago_id
INNER JOIN usuarios u ON u.id = v.usuario_id
WHERE v.id = $1;

-- name: ListVentasPaginated :many
SELECT v.*, mp.nombre AS metodo_pago_nombre, mp.codigo AS metodo_pago_codigo,
       u.nombre_completo AS usuario_nombre,
       (SELECT COUNT(*) FROM venta_items WHERE venta_id = v.id) AS items_count
FROM ventas v
INNER JOIN metodos_pago mp ON mp.id = v.metodo_pago_id
INNER JOIN usuarios u ON u.id = v.usuario_id
WHERE
    (@sede_id::bigint = 0 OR v.sede_id = @sede_id)
    AND (@usuario_id::bigint = 0 OR v.usuario_id = @usuario_id)
    AND (@metodo_pago_id::bigint = 0 OR v.metodo_pago_id = @metodo_pago_id)
    AND (@fecha_desde::timestamptz IS NULL OR v.created_at >= @fecha_desde)
    AND (@fecha_hasta::timestamptz IS NULL OR v.created_at <= @fecha_hasta)
    AND (@total_min::numeric = 0 OR v.total >= @total_min)
    AND (@total_max::numeric = 0 OR v.total <= @total_max)
    AND (NOT @con_descuento::bool OR v.descuento_monto > 0)
ORDER BY v.created_at DESC, v.id DESC
LIMIT $1 OFFSET $2;

-- name: CountVentas :one
-- Mismos filtros que ListVentasPaginated (sin joins ni columnas extra) —
-- el spec ya los trajo pareados; verificado que no hay drift.
SELECT COUNT(*) FROM ventas v
WHERE
    (@sede_id::bigint = 0 OR v.sede_id = @sede_id)
    AND (@usuario_id::bigint = 0 OR v.usuario_id = @usuario_id)
    AND (@metodo_pago_id::bigint = 0 OR v.metodo_pago_id = @metodo_pago_id)
    AND (@fecha_desde::timestamptz IS NULL OR v.created_at >= @fecha_desde)
    AND (@fecha_hasta::timestamptz IS NULL OR v.created_at <= @fecha_hasta)
    AND (@total_min::numeric = 0 OR v.total >= @total_min)
    AND (@total_max::numeric = 0 OR v.total <= @total_max)
    AND (NOT @con_descuento::bool OR v.descuento_monto > 0);

-- name: UpdateVentaObservaciones :one
UPDATE ventas SET observaciones = $2 WHERE id = $1 RETURNING *;

-- name: GetResumenVentasHoy :one
WITH ventas_hoy AS (
    SELECT * FROM ventas v_raw
    WHERE v_raw.sede_id = @sede_id::bigint
      AND v_raw.created_at >= @dia_inicio::timestamptz
      AND v_raw.created_at < @dia_fin::timestamptz
),
totales_por_metodo AS (
    SELECT mp.codigo, COALESCE(SUM(v.total), 0)::numeric AS total_metodo
    FROM metodos_pago mp
    LEFT JOIN ventas_hoy v ON v.metodo_pago_id = mp.id
    GROUP BY mp.codigo
)
SELECT
    COALESCE((SELECT SUM(total_metodo) FROM totales_por_metodo WHERE codigo = 'efectivo'), 0)::numeric AS total_efectivo,
    COALESCE((SELECT SUM(total_metodo) FROM totales_por_metodo WHERE codigo = 'nequi'), 0)::numeric AS total_nequi,
    COALESCE((SELECT SUM(total_metodo) FROM totales_por_metodo WHERE codigo = 'daviplata'), 0)::numeric AS total_daviplata,
    COALESCE((SELECT SUM(total_metodo) FROM totales_por_metodo WHERE codigo = 'transferencia'), 0)::numeric AS total_transferencia,
    COALESCE((SELECT SUM(total_metodo) FROM totales_por_metodo WHERE codigo NOT IN ('efectivo','nequi','daviplata','transferencia')), 0)::numeric AS total_otros,
    (SELECT COUNT(*) FROM ventas_hoy)::bigint AS ventas_count,
    COALESCE((SELECT SUM(total) FROM ventas_hoy), 0)::numeric AS total_dia,
    COALESCE((SELECT SUM(descuento_monto) FROM ventas_hoy), 0)::numeric AS descuento_total;

-- name: GetVentasPorVendedoraHoy :many
SELECT u.id AS usuario_id, u.nombre_completo, COUNT(v.id)::bigint AS ventas_count, COALESCE(SUM(v.total), 0)::numeric AS total
FROM ventas v
INNER JOIN usuarios u ON u.id = v.usuario_id
WHERE v.sede_id = @sede_id
  AND v.created_at >= @dia_inicio
  AND v.created_at < @dia_fin
GROUP BY u.id, u.nombre_completo
ORDER BY total DESC;

-- name: GetTopFraganciasHoy :many
SELECT f.id, f.nombre_comercial,
       COALESCE(SUM(vi.gramos_fragancia * vi.cantidad), 0)::numeric AS gramos_vendidos,
       COALESCE(SUM(vi.subtotal), 0)::numeric AS monto_vendido
FROM venta_items vi
INNER JOIN ventas v ON v.id = vi.venta_id
INNER JOIN fragancias f ON f.id = vi.fragancia_id
WHERE v.sede_id = @sede_id
  AND v.created_at >= @dia_inicio
  AND v.created_at < @dia_fin
  AND vi.fragancia_id IS NOT NULL
GROUP BY f.id, f.nombre_comercial
ORDER BY monto_vendido DESC
LIMIT 10;
