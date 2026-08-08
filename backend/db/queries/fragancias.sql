-- name: ListFraganciasPaginated :many
SELECT f.*,
  COALESCE(SUM(CASE WHEN sa.ubicacion = 'vitrina' THEN sa.cantidad ELSE 0 END), 0)::numeric AS stock_vitrina,
  COALESCE(SUM(CASE WHEN sa.ubicacion = 'bodega' THEN sa.cantidad ELSE 0 END), 0)::numeric AS stock_bodega
FROM fragancias f
LEFT JOIN stock_actual sa ON sa.tipo_item = 'fragancia' AND sa.item_id = f.id
WHERE
  (@include_deleted::bool OR f.deleted_at IS NULL)
  AND (@sede_id::bigint = 0 OR f.sede_id = @sede_id)
  AND (@genero::text = '' OR f.genero::text = @genero)
  AND (@numero_genero::int = 0 OR f.numero_genero = @numero_genero)
  AND (@activo::text = 'all' OR (@activo::text = 'true' AND f.activo = true) OR (@activo::text = 'false' AND f.activo = false))
  AND (@q::text = '' OR f.nombre_comercial ILIKE '%' || @q || '%' OR f.nombre_alternativo ILIKE '%' || @q || '%')
GROUP BY f.id
HAVING
  NOT @stock_bajo::bool OR (
    COALESCE(SUM(CASE WHEN sa.ubicacion = 'vitrina' THEN sa.cantidad ELSE 0 END), 0)
    + COALESCE(SUM(CASE WHEN sa.ubicacion = 'bodega' THEN sa.cantidad ELSE 0 END), 0)
    < f.gramos_minimo
  )
ORDER BY
  CASE WHEN @sort_col::text = 'nombre_comercial' AND @sort_dir::text = 'asc' THEN f.nombre_comercial END ASC,
  CASE WHEN @sort_col::text = 'nombre_comercial' AND @sort_dir::text = 'desc' THEN f.nombre_comercial END DESC,
  CASE WHEN @sort_col::text = 'created_at' AND @sort_dir::text = 'asc' THEN f.created_at END ASC,
  CASE WHEN @sort_col::text = 'created_at' AND @sort_dir::text = 'desc' THEN f.created_at END DESC,
  f.id ASC
LIMIT $1 OFFSET $2;

-- name: CountFragancias :one
-- Mirrors ListFraganciasPaginated's GROUP BY/HAVING (stock_bajo needs it),
-- otherwise meta.total would drift from the actually-returned items.
SELECT COUNT(*) FROM (
  SELECT f.id
  FROM fragancias f
  LEFT JOIN stock_actual sa ON sa.tipo_item = 'fragancia' AND sa.item_id = f.id
  WHERE
    (@include_deleted::bool OR f.deleted_at IS NULL)
    AND (@sede_id::bigint = 0 OR f.sede_id = @sede_id)
    AND (@genero::text = '' OR f.genero::text = @genero)
    AND (@numero_genero::int = 0 OR f.numero_genero = @numero_genero)
    AND (@activo::text = 'all' OR (@activo::text = 'true' AND f.activo = true) OR (@activo::text = 'false' AND f.activo = false))
    AND (@q::text = '' OR f.nombre_comercial ILIKE '%' || @q || '%' OR f.nombre_alternativo ILIKE '%' || @q || '%')
  GROUP BY f.id
  HAVING
    NOT @stock_bajo::bool OR (
      COALESCE(SUM(CASE WHEN sa.ubicacion = 'vitrina' THEN sa.cantidad ELSE 0 END), 0)
      + COALESCE(SUM(CASE WHEN sa.ubicacion = 'bodega' THEN sa.cantidad ELSE 0 END), 0)
      < f.gramos_minimo
    )
) sub;

-- name: GetFraganciaByID :one
SELECT f.*,
  COALESCE(SUM(CASE WHEN sa.ubicacion = 'vitrina' THEN sa.cantidad ELSE 0 END), 0)::numeric AS stock_vitrina,
  COALESCE(SUM(CASE WHEN sa.ubicacion = 'bodega' THEN sa.cantidad ELSE 0 END), 0)::numeric AS stock_bodega
FROM fragancias f
LEFT JOIN stock_actual sa ON sa.tipo_item = 'fragancia' AND sa.item_id = f.id
WHERE f.id = $1 AND f.deleted_at IS NULL
GROUP BY f.id;

-- name: GetFraganciaByIDIncludingDeleted :one
SELECT * FROM fragancias WHERE id = $1;

-- name: InsertFragancia :one
INSERT INTO fragancias (sede_id, nombre_comercial, nombre_alternativo, genero, gramos_minimo, numero_genero)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateFragancia :one
UPDATE fragancias SET
  nombre_comercial = COALESCE(sqlc.narg('nombre_comercial'), nombre_comercial),
  nombre_alternativo = COALESCE(sqlc.narg('nombre_alternativo'), nombre_alternativo),
  genero = COALESCE(sqlc.narg('genero'), genero),
  gramos_minimo = COALESCE(sqlc.narg('gramos_minimo'), gramos_minimo),
  numero_genero = COALESCE(sqlc.narg('numero_genero'), numero_genero),
  updated_at = NOW()
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteFragancia :exec
UPDATE fragancias SET deleted_at = NOW(), activo = false, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: RestoreFragancia :one
UPDATE fragancias SET deleted_at = NULL, activo = true, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NOT NULL
RETURNING *;

-- name: ExistsFraganciaNombreComercial :one
SELECT EXISTS(
  SELECT 1 FROM fragancias
  WHERE sede_id = @sede_id AND LOWER(nombre_comercial) = LOWER(@nombre_comercial) AND deleted_at IS NULL
    AND (@exclude_id::bigint = 0 OR id != @exclude_id)
);

-- name: ExistsFraganciaNumeroGenero :one
SELECT EXISTS(
  SELECT 1 FROM fragancias
  WHERE sede_id = @sede_id AND genero = @genero AND numero_genero = @numero_genero AND deleted_at IS NULL
    AND (@exclude_id::bigint = 0 OR id != @exclude_id)
);

-- name: NextNumeroGeneroFragancia :one
SELECT COALESCE(MAX(numero_genero), 0) + 1 AS siguiente
FROM fragancias
WHERE sede_id = $1 AND genero = $2 AND deleted_at IS NULL;
