-- name: ListMetodosPagoPaginated :many
SELECT * FROM metodos_pago
WHERE
  (@include_deleted::bool OR deleted_at IS NULL)
  AND (@activo::text = 'all' OR (@activo::text = 'true' AND activo = true) OR (@activo::text = 'false' AND activo = false))
  AND (@q::text = '' OR nombre ILIKE '%' || @q || '%' OR codigo ILIKE '%' || @q || '%')
ORDER BY
  CASE WHEN @sort_col::text = 'nombre' AND @sort_dir::text = 'asc' THEN nombre END ASC,
  CASE WHEN @sort_col::text = 'nombre' AND @sort_dir::text = 'desc' THEN nombre END DESC,
  CASE WHEN @sort_col::text = 'created_at' AND @sort_dir::text = 'asc' THEN created_at END ASC,
  CASE WHEN @sort_col::text = 'created_at' AND @sort_dir::text = 'desc' THEN created_at END DESC,
  id ASC
LIMIT $1 OFFSET $2;

-- name: CountMetodosPago :one
-- No HAVING/aggregate here, so a plain COUNT can't drift from the list.
SELECT COUNT(*) FROM metodos_pago
WHERE
  (@include_deleted::bool OR deleted_at IS NULL)
  AND (@activo::text = 'all' OR (@activo::text = 'true' AND activo = true) OR (@activo::text = 'false' AND activo = false))
  AND (@q::text = '' OR nombre ILIKE '%' || @q || '%' OR codigo ILIKE '%' || @q || '%');

-- name: GetMetodoPagoByID :one
SELECT * FROM metodos_pago WHERE id = $1 AND deleted_at IS NULL;

-- name: GetMetodoPagoByIDIncludingDeleted :one
SELECT * FROM metodos_pago WHERE id = $1;

-- name: InsertMetodoPago :one
INSERT INTO metodos_pago (nombre, codigo)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateMetodoPago :one
UPDATE metodos_pago SET
  nombre = COALESCE(sqlc.narg('nombre'), nombre),
  codigo = COALESCE(sqlc.narg('codigo'), codigo),
  updated_at = NOW()
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteMetodoPago :exec
UPDATE metodos_pago SET deleted_at = NOW(), activo = false, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: HardDeleteMetodoPago :exec
DELETE FROM metodos_pago WHERE id = $1;

-- name: CountVentasByMetodoPago :one
SELECT COUNT(*) FROM ventas WHERE metodo_pago_id = $1;

-- name: ExistsMetodoPagoCodigo :one
SELECT EXISTS(
  SELECT 1 FROM metodos_pago
  WHERE LOWER(codigo) = LOWER(@codigo) AND deleted_at IS NULL
    AND (@exclude_id::bigint = 0 OR id != @exclude_id)
);

-- name: ExistsMetodoPagoNombre :one
SELECT EXISTS(
  SELECT 1 FROM metodos_pago
  WHERE LOWER(nombre) = LOWER(@nombre) AND deleted_at IS NULL
    AND (@exclude_id::bigint = 0 OR id != @exclude_id)
);
