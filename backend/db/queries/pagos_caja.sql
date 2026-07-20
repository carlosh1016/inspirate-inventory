-- name: InsertPagoCaja :one
INSERT INTO pagos_caja (cuadre_caja_id, usuario_id, concepto, monto)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPagosByCuadre :many
SELECT p.*, u.nombre_completo AS usuario_nombre
FROM pagos_caja p
INNER JOIN usuarios u ON u.id = p.usuario_id
WHERE p.cuadre_caja_id = $1
ORDER BY p.created_at ASC;

-- name: GetPagoByID :one
SELECT * FROM pagos_caja WHERE id = $1;

-- name: DeletePagoCaja :exec
DELETE FROM pagos_caja WHERE id = $1;

-- name: GetTotalPagosByCuadre :one
SELECT COALESCE(SUM(monto), 0)::numeric AS total FROM pagos_caja WHERE cuadre_caja_id = $1;
