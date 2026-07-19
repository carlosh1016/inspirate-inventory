-- name: ListProductosPaginated :many
SELECT p.*,
  COALESCE(SUM(CASE WHEN sa.ubicacion = 'vitrina' THEN sa.cantidad ELSE 0 END), 0)::numeric AS stock_vitrina,
  COALESCE(SUM(CASE WHEN sa.ubicacion = 'bodega' THEN sa.cantidad ELSE 0 END), 0)::numeric AS stock_bodega
FROM productos p
LEFT JOIN stock_actual sa ON sa.tipo_item = 'producto' AND sa.item_id = p.id
WHERE
  (@include_deleted::bool OR p.deleted_at IS NULL)
  AND (@sede_id::bigint = 0 OR p.sede_id = @sede_id)
  AND (@categoria::text = '' OR p.categoria::text = @categoria)
  AND (@activo::text = 'all' OR (@activo::text = 'true' AND p.activo = true) OR (@activo::text = 'false' AND p.activo = false))
  AND (@q::text = '' OR p.nombre ILIKE '%' || @q || '%')
GROUP BY p.id
HAVING
  NOT @stock_bajo::bool OR (
    COALESCE(SUM(CASE WHEN sa.ubicacion = 'vitrina' THEN sa.cantidad ELSE 0 END), 0)
    + COALESCE(SUM(CASE WHEN sa.ubicacion = 'bodega' THEN sa.cantidad ELSE 0 END), 0)
    < p.stock_minimo
  )
ORDER BY
  CASE WHEN @sort_col::text = 'nombre' AND @sort_dir::text = 'asc' THEN p.nombre END ASC,
  CASE WHEN @sort_col::text = 'nombre' AND @sort_dir::text = 'desc' THEN p.nombre END DESC,
  CASE WHEN @sort_col::text = 'created_at' AND @sort_dir::text = 'asc' THEN p.created_at END ASC,
  CASE WHEN @sort_col::text = 'created_at' AND @sort_dir::text = 'desc' THEN p.created_at END DESC,
  p.id ASC
LIMIT $1 OFFSET $2;

-- name: CountProductos :one
-- Mirrors ListProductosPaginated's GROUP BY/HAVING (stock_bajo needs it),
-- otherwise meta.total would drift from the actually-returned items — the
-- same bug fixed in fragancias' CountFragancias.
SELECT COUNT(*) FROM (
  SELECT p.id
  FROM productos p
  LEFT JOIN stock_actual sa ON sa.tipo_item = 'producto' AND sa.item_id = p.id
  WHERE
    (@include_deleted::bool OR p.deleted_at IS NULL)
    AND (@sede_id::bigint = 0 OR p.sede_id = @sede_id)
    AND (@categoria::text = '' OR p.categoria::text = @categoria)
    AND (@activo::text = 'all' OR (@activo::text = 'true' AND p.activo = true) OR (@activo::text = 'false' AND p.activo = false))
    AND (@q::text = '' OR p.nombre ILIKE '%' || @q || '%')
  GROUP BY p.id
  HAVING
    NOT @stock_bajo::bool OR (
      COALESCE(SUM(CASE WHEN sa.ubicacion = 'vitrina' THEN sa.cantidad ELSE 0 END), 0)
      + COALESCE(SUM(CASE WHEN sa.ubicacion = 'bodega' THEN sa.cantidad ELSE 0 END), 0)
      < p.stock_minimo
    )
) sub;

-- name: GetProductoByID :one
SELECT p.*,
  COALESCE(SUM(CASE WHEN sa.ubicacion = 'vitrina' THEN sa.cantidad ELSE 0 END), 0)::numeric AS stock_vitrina,
  COALESCE(SUM(CASE WHEN sa.ubicacion = 'bodega' THEN sa.cantidad ELSE 0 END), 0)::numeric AS stock_bodega
FROM productos p
LEFT JOIN stock_actual sa ON sa.tipo_item = 'producto' AND sa.item_id = p.id
WHERE p.id = $1 AND p.deleted_at IS NULL
GROUP BY p.id;

-- name: GetProductoByIDIncludingDeleted :one
SELECT * FROM productos WHERE id = $1;

-- name: InsertProducto :one
INSERT INTO productos (sede_id, nombre, categoria, precio, stock_minimo)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateProducto :one
UPDATE productos SET
  nombre = COALESCE(sqlc.narg('nombre'), nombre),
  categoria = COALESCE(sqlc.narg('categoria'), categoria),
  precio = COALESCE(sqlc.narg('precio'), precio),
  stock_minimo = COALESCE(sqlc.narg('stock_minimo'), stock_minimo),
  updated_at = NOW()
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteProducto :exec
UPDATE productos SET deleted_at = NOW(), activo = false, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: ExistsProductoNombreCategoria :one
SELECT EXISTS(
  SELECT 1 FROM productos
  WHERE sede_id = @sede_id AND LOWER(nombre) = LOWER(@nombre) AND categoria = @categoria::categoria_producto_enum AND deleted_at IS NULL
    AND (@exclude_id::bigint = 0 OR id != @exclude_id)
);
