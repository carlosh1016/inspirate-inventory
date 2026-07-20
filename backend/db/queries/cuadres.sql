-- name: GetCuadreByID :one
SELECT c.*, u.nombre_completo AS cerrado_por_nombre
FROM cuadres_caja c
LEFT JOIN usuarios u ON u.id = c.cerrado_por_usuario_id
WHERE c.id = $1;

-- name: GetCuadreBySedeFecha :one
SELECT c.*, u.nombre_completo AS cerrado_por_nombre
FROM cuadres_caja c
LEFT JOIN usuarios u ON u.id = c.cerrado_por_usuario_id
WHERE c.sede_id = $1 AND c.fecha = $2;

-- name: GetCuadreAbiertoAnterior :one
-- Most recent still-open cuadre strictly before fecha, for the "opened a
-- new day while yesterday's is still open" soft warning.
SELECT c.*
FROM cuadres_caja c
WHERE c.sede_id = $1 AND c.estado = 'abierto' AND c.fecha < $2
ORDER BY c.fecha DESC
LIMIT 1;

-- name: ListCuadresPaginated :many
SELECT c.*, u.nombre_completo AS cerrado_por_nombre
FROM cuadres_caja c
LEFT JOIN usuarios u ON u.id = c.cerrado_por_usuario_id
WHERE
    (@sede_id::bigint = 0 OR c.sede_id = @sede_id)
    AND (@estado::text = '' OR c.estado::text = @estado)
    AND (@fecha_desde::date IS NULL OR c.fecha >= @fecha_desde)
    AND (@fecha_hasta::date IS NULL OR c.fecha <= @fecha_hasta)
ORDER BY c.fecha DESC
LIMIT $1 OFFSET $2;

-- name: CountCuadres :one
-- Mismos filtros que ListCuadresPaginated (sin joins ni columnas extra).
SELECT COUNT(*) FROM cuadres_caja c
WHERE
    (@sede_id::bigint = 0 OR c.sede_id = @sede_id)
    AND (@estado::text = '' OR c.estado::text = @estado)
    AND (@fecha_desde::date IS NULL OR c.fecha >= @fecha_desde)
    AND (@fecha_hasta::date IS NULL OR c.fecha <= @fecha_hasta);

-- name: InsertCuadre :one
INSERT INTO cuadres_caja (sede_id, fecha, estado, fondo_base)
VALUES ($1, $2, 'abierto', $3)
RETURNING *;

-- name: UpdateCuadreCerrar :one
UPDATE cuadres_caja SET
    estado = 'cerrado',
    total_efectivo = $2,
    total_nequi = $3,
    total_daviplata = $4,
    total_transferencia = $5,
    total_otros = $6,
    total_pagos = $7,
    total_consignaciones = $8,
    valor_turno = $9,
    saldo_calculado = $10,
    observaciones = $11,
    cerrado_por_usuario_id = $12,
    cerrado_at = NOW(),
    updated_at = NOW()
WHERE id = $1 AND estado = 'abierto'
RETURNING *;

-- name: ExistsCuadreCerradoBySedeFecha :one
SELECT EXISTS(
    SELECT 1 FROM cuadres_caja
    WHERE sede_id = $1 AND fecha = $2 AND estado = 'cerrado'
);

-- name: GetTotalesPorMetodoEnFecha :one
SELECT
    COALESCE(SUM(CASE WHEN mp.codigo = 'efectivo' THEN v.total ELSE 0 END), 0)::numeric AS total_efectivo,
    COALESCE(SUM(CASE WHEN mp.codigo = 'nequi' THEN v.total ELSE 0 END), 0)::numeric AS total_nequi,
    COALESCE(SUM(CASE WHEN mp.codigo = 'daviplata' THEN v.total ELSE 0 END), 0)::numeric AS total_daviplata,
    COALESCE(SUM(CASE WHEN mp.codigo = 'transferencia' THEN v.total ELSE 0 END), 0)::numeric AS total_transferencia,
    COALESCE(SUM(CASE WHEN mp.codigo NOT IN ('efectivo','nequi','daviplata','transferencia') THEN v.total ELSE 0 END), 0)::numeric AS total_otros,
    COALESCE(SUM(v.total), 0)::numeric AS total_dia,
    COUNT(*)::bigint AS ventas_count
FROM ventas v
INNER JOIN metodos_pago mp ON mp.id = v.metodo_pago_id
WHERE v.sede_id = @sede_id
  AND v.created_at >= @dia_inicio
  AND v.created_at < @dia_fin;
