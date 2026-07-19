-- name: ListModelosEnvasePaginated :many
SELECT m.*,
  COUNT(v.id) FILTER (WHERE v.deleted_at IS NULL) AS variantes_activas
FROM modelos_envase m
LEFT JOIN variantes_envase v ON v.modelo_envase_id = m.id
WHERE
  (@include_deleted::bool OR m.deleted_at IS NULL)
  AND (@tipo::text = '' OR m.tipo ILIKE @tipo)
  AND (@activo::text = 'all' OR (@activo::text = 'true' AND m.activo = true) OR (@activo::text = 'false' AND m.activo = false))
  AND (@q::text = '' OR m.tipo ILIKE '%' || @q || '%')
GROUP BY m.id
ORDER BY
  CASE WHEN @sort_col::text = 'tipo' AND @sort_dir::text = 'asc' THEN m.tipo END ASC,
  CASE WHEN @sort_col::text = 'tipo' AND @sort_dir::text = 'desc' THEN m.tipo END DESC,
  CASE WHEN @sort_col::text = 'created_at' AND @sort_dir::text = 'asc' THEN m.created_at END ASC,
  CASE WHEN @sort_col::text = 'created_at' AND @sort_dir::text = 'desc' THEN m.created_at END DESC,
  m.id ASC
LIMIT $1 OFFSET $2;

-- name: CountModelosEnvase :one
-- No HAVING here: unlike fragancias' stock_bajo, nothing filters on the
-- variantes_activas aggregate, so a plain COUNT can't drift from the list.
SELECT COUNT(*) FROM modelos_envase m
WHERE
  (@include_deleted::bool OR m.deleted_at IS NULL)
  AND (@tipo::text = '' OR m.tipo ILIKE @tipo)
  AND (@activo::text = 'all' OR (@activo::text = 'true' AND m.activo = true) OR (@activo::text = 'false' AND m.activo = false))
  AND (@q::text = '' OR m.tipo ILIKE '%' || @q || '%');

-- name: GetModeloEnvaseByID :one
SELECT m.*,
  COUNT(v.id) FILTER (WHERE v.deleted_at IS NULL) AS variantes_activas
FROM modelos_envase m
LEFT JOIN variantes_envase v ON v.modelo_envase_id = m.id
WHERE m.id = $1 AND m.deleted_at IS NULL
GROUP BY m.id;

-- name: GetModeloEnvaseByIDIncludingDeleted :one
SELECT * FROM modelos_envase WHERE id = $1;

-- name: InsertModeloEnvase :one
INSERT INTO modelos_envase (tipo, tamano_oz, equiv_gramos, precio_solo, precio_con_fragancia, precio_recarga)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateModeloEnvase :one
UPDATE modelos_envase SET
  tipo = COALESCE(sqlc.narg('tipo'), tipo),
  tamano_oz = COALESCE(sqlc.narg('tamano_oz'), tamano_oz),
  equiv_gramos = COALESCE(sqlc.narg('equiv_gramos'), equiv_gramos),
  precio_solo = COALESCE(sqlc.narg('precio_solo'), precio_solo),
  precio_con_fragancia = COALESCE(sqlc.narg('precio_con_fragancia'), precio_con_fragancia),
  precio_recarga = COALESCE(sqlc.narg('precio_recarga'), precio_recarga),
  updated_at = NOW()
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteModeloEnvase :exec
UPDATE modelos_envase SET deleted_at = NOW(), activo = false, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: CountVariantesActivasByModelo :one
SELECT COUNT(*) FROM variantes_envase
WHERE modelo_envase_id = $1 AND deleted_at IS NULL;

-- name: ExistsModeloEnvaseTipoTamano :one
SELECT EXISTS(
  SELECT 1 FROM modelos_envase
  WHERE LOWER(tipo) = LOWER(@tipo) AND tamano_oz = @tamano_oz AND deleted_at IS NULL
    AND (@exclude_id::bigint = 0 OR id != @exclude_id)
);
