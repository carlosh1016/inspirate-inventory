-- name: InsertSesion :one
INSERT INTO sesiones_laborales (sede_id, usuario_id, entrada_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSesionAbiertaPorUsuario :one
SELECT * FROM sesiones_laborales
WHERE usuario_id = $1 AND salida_at IS NULL
LIMIT 1;

-- name: CerrarSesion :one
UPDATE sesiones_laborales SET
    salida_at = $2,
    horas_trabajadas = $2 - entrada_at,
    updated_at = NOW()
WHERE id = $1 AND salida_at IS NULL
RETURNING *;

-- name: GetSesionByID :one
SELECT s.*, u.nombre_completo AS usuario_nombre
FROM sesiones_laborales s
INNER JOIN usuarios u ON u.id = s.usuario_id
WHERE s.id = $1;

-- name: ListSesiones :many
SELECT s.*, u.nombre_completo AS usuario_nombre
FROM sesiones_laborales s
INNER JOIN usuarios u ON u.id = s.usuario_id
WHERE
    (@sede_id::bigint = 0 OR s.sede_id = @sede_id)
    AND (@usuario_id::bigint = 0 OR s.usuario_id = @usuario_id)
    AND (@fecha_desde::timestamptz IS NULL OR s.entrada_at >= @fecha_desde)
    AND (@fecha_hasta::timestamptz IS NULL OR s.entrada_at <= @fecha_hasta)
    AND (NOT @abiertas::bool OR s.salida_at IS NULL)
ORDER BY s.entrada_at DESC
LIMIT $1 OFFSET $2;

-- name: CountSesiones :one
-- Mismos filtros que ListSesiones (sin joins ni columnas extra).
SELECT COUNT(*) FROM sesiones_laborales s
WHERE
    (@sede_id::bigint = 0 OR s.sede_id = @sede_id)
    AND (@usuario_id::bigint = 0 OR s.usuario_id = @usuario_id)
    AND (@fecha_desde::timestamptz IS NULL OR s.entrada_at >= @fecha_desde)
    AND (@fecha_hasta::timestamptz IS NULL OR s.entrada_at <= @fecha_hasta)
    AND (NOT @abiertas::bool OR s.salida_at IS NULL);

-- name: UpdateSesionManual :one
UPDATE sesiones_laborales SET
    entrada_at = COALESCE(sqlc.narg('entrada_at'), entrada_at),
    salida_at = COALESCE(sqlc.narg('salida_at'), salida_at),
    horas_trabajadas = CASE
        WHEN COALESCE(sqlc.narg('salida_at'), salida_at) IS NOT NULL
            THEN COALESCE(sqlc.narg('salida_at'), salida_at) - COALESCE(sqlc.narg('entrada_at'), entrada_at)
        ELSE NULL
    END,
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: GetResumenSesiones :many
SELECT
    u.id AS usuario_id,
    u.nombre_completo,
    COUNT(*)::bigint AS sesiones_count,
    COUNT(DISTINCT DATE(s.entrada_at AT TIME ZONE 'America/Bogota'))::bigint AS dias_trabajados,
    COALESCE(SUM(s.horas_trabajadas), INTERVAL '0')::interval AS total_horas
FROM sesiones_laborales s
INNER JOIN usuarios u ON u.id = s.usuario_id
WHERE s.entrada_at >= @fecha_desde
  AND s.entrada_at <= @fecha_hasta
  AND (@usuario_id::bigint = 0 OR s.usuario_id = @usuario_id)
  AND s.salida_at IS NOT NULL
GROUP BY u.id, u.nombre_completo
ORDER BY u.nombre_completo ASC;
