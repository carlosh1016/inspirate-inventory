-- name: InsertAuditoria :exec
INSERT INTO auditoria (usuario_id, accion, tabla_afectada, registro_id, datos_antes, datos_despues, ip, user_agent)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: ListAuditoriaPaginated :many
-- usuario_nombre viene de un LEFT JOIN (usuario_id puede ser NULL, p.ej. en
-- login_failed sin correo válido) → sqlc lo infiere nullable (pgtype.Text).
SELECT a.*, u.nombre_completo AS usuario_nombre
FROM auditoria a
LEFT JOIN usuarios u ON u.id = a.usuario_id
WHERE
    (@usuario_id::bigint = 0 OR a.usuario_id = @usuario_id)
    AND (@accion::text = '' OR a.accion = @accion)
    AND (@tabla_afectada::text = '' OR a.tabla_afectada = @tabla_afectada)
    AND (@fecha_desde::timestamptz IS NULL OR a.created_at >= @fecha_desde)
    AND (@fecha_hasta::timestamptz IS NULL OR a.created_at <= @fecha_hasta)
ORDER BY a.created_at DESC, a.id DESC
LIMIT $1 OFFSET $2;

-- name: CountAuditoria :one
-- Espejo exacto de los filtros de ListAuditoriaPaginated.
SELECT COUNT(*) FROM auditoria a
WHERE
    (@usuario_id::bigint = 0 OR a.usuario_id = @usuario_id)
    AND (@accion::text = '' OR a.accion = @accion)
    AND (@tabla_afectada::text = '' OR a.tabla_afectada = @tabla_afectada)
    AND (@fecha_desde::timestamptz IS NULL OR a.created_at >= @fecha_desde)
    AND (@fecha_hasta::timestamptz IS NULL OR a.created_at <= @fecha_hasta);

-- name: GetAuditoriaByID :one
SELECT a.*, u.nombre_completo AS usuario_nombre
FROM auditoria a
LEFT JOIN usuarios u ON u.id = a.usuario_id
WHERE a.id = $1;

-- name: GetAccionesDistintas :many
SELECT DISTINCT accion FROM auditoria ORDER BY accion ASC;
