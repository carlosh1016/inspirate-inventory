-- name: InsertConsignacion :one
INSERT INTO consignaciones (cuadre_caja_id, usuario_id, monto, banco, referencia)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetConsignacionesByCuadre :many
SELECT c.*, u.nombre_completo AS usuario_nombre
FROM consignaciones c
INNER JOIN usuarios u ON u.id = c.usuario_id
WHERE c.cuadre_caja_id = $1
ORDER BY c.created_at ASC;

-- name: GetConsignacionByID :one
SELECT * FROM consignaciones WHERE id = $1;

-- name: DeleteConsignacion :exec
DELETE FROM consignaciones WHERE id = $1;

-- name: GetTotalConsignacionesByCuadre :one
SELECT COALESCE(SUM(monto), 0)::numeric AS total FROM consignaciones WHERE cuadre_caja_id = $1;
