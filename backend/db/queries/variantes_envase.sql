-- name: ListVariantesEnvasePaginated :many
SELECT v.*,
  COALESCE(SUM(CASE WHEN sa.ubicacion = 'vitrina' THEN sa.cantidad ELSE 0 END), 0)::numeric AS stock_vitrina,
  COALESCE(SUM(CASE WHEN sa.ubicacion = 'bodega' THEN sa.cantidad ELSE 0 END), 0)::numeric AS stock_bodega
FROM variantes_envase v
LEFT JOIN stock_actual sa ON sa.tipo_item = 'variante_envase' AND sa.item_id = v.id
WHERE
  (@include_deleted::bool OR v.deleted_at IS NULL)
  AND (@sede_id::bigint = 0 OR v.sede_id = @sede_id)
  AND (@modelo_envase_id::bigint = 0 OR v.modelo_envase_id = @modelo_envase_id)
  AND (@activo::text = 'all' OR (@activo::text = 'true' AND v.activo = true) OR (@activo::text = 'false' AND v.activo = false))
  AND (@q::text = '' OR v.color ILIKE '%' || @q || '%')
GROUP BY v.id
HAVING
  NOT @stock_bajo::bool OR (
    COALESCE(SUM(CASE WHEN sa.ubicacion = 'vitrina' THEN sa.cantidad ELSE 0 END), 0)
    + COALESCE(SUM(CASE WHEN sa.ubicacion = 'bodega' THEN sa.cantidad ELSE 0 END), 0)
    < v.stock_minimo
  )
ORDER BY
  CASE WHEN @sort_col::text = 'color' AND @sort_dir::text = 'asc' THEN v.color END ASC,
  CASE WHEN @sort_col::text = 'color' AND @sort_dir::text = 'desc' THEN v.color END DESC,
  CASE WHEN @sort_col::text = 'created_at' AND @sort_dir::text = 'asc' THEN v.created_at END ASC,
  CASE WHEN @sort_col::text = 'created_at' AND @sort_dir::text = 'desc' THEN v.created_at END DESC,
  v.id ASC
LIMIT $1 OFFSET $2;

-- name: CountVariantesEnvase :one
-- Mirrors ListVariantesEnvasePaginated's GROUP BY/HAVING (stock_bajo needs
-- it), otherwise meta.total would drift from the actually-returned items —
-- the same bug fixed in fragancias' CountFragancias.
SELECT COUNT(*) FROM (
  SELECT v.id
  FROM variantes_envase v
  LEFT JOIN stock_actual sa ON sa.tipo_item = 'variante_envase' AND sa.item_id = v.id
  WHERE
    (@include_deleted::bool OR v.deleted_at IS NULL)
    AND (@sede_id::bigint = 0 OR v.sede_id = @sede_id)
    AND (@modelo_envase_id::bigint = 0 OR v.modelo_envase_id = @modelo_envase_id)
    AND (@activo::text = 'all' OR (@activo::text = 'true' AND v.activo = true) OR (@activo::text = 'false' AND v.activo = false))
    AND (@q::text = '' OR v.color ILIKE '%' || @q || '%')
  GROUP BY v.id
  HAVING
    NOT @stock_bajo::bool OR (
      COALESCE(SUM(CASE WHEN sa.ubicacion = 'vitrina' THEN sa.cantidad ELSE 0 END), 0)
      + COALESCE(SUM(CASE WHEN sa.ubicacion = 'bodega' THEN sa.cantidad ELSE 0 END), 0)
      < v.stock_minimo
    )
) sub;

-- name: GetVarianteEnvaseByID :one
SELECT v.*,
  COALESCE(SUM(CASE WHEN sa.ubicacion = 'vitrina' THEN sa.cantidad ELSE 0 END), 0)::numeric AS stock_vitrina,
  COALESCE(SUM(CASE WHEN sa.ubicacion = 'bodega' THEN sa.cantidad ELSE 0 END), 0)::numeric AS stock_bodega
FROM variantes_envase v
LEFT JOIN stock_actual sa ON sa.tipo_item = 'variante_envase' AND sa.item_id = v.id
WHERE v.id = $1 AND v.deleted_at IS NULL
GROUP BY v.id;

-- name: GetVarianteEnvaseByIDIncludingDeleted :one
SELECT * FROM variantes_envase WHERE id = $1;

-- name: InsertVarianteEnvase :one
INSERT INTO variantes_envase (modelo_envase_id, sede_id, color, stock_minimo)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateVarianteEnvase :one
UPDATE variantes_envase SET
  color = COALESCE(sqlc.narg('color'), color),
  stock_minimo = COALESCE(sqlc.narg('stock_minimo'), stock_minimo),
  updated_at = NOW()
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteVarianteEnvase :exec
UPDATE variantes_envase SET deleted_at = NOW(), activo = false, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: ExistsVarianteEnvaseColor :one
SELECT EXISTS(
  SELECT 1 FROM variantes_envase
  WHERE modelo_envase_id = @modelo_envase_id AND LOWER(color) = LOWER(@color) AND deleted_at IS NULL
    AND (@exclude_id::bigint = 0 OR id != @exclude_id)
);
