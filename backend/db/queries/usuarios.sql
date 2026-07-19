-- name: GetUsuarioByCorreo :one
SELECT * FROM usuarios
WHERE LOWER(correo) = LOWER($1) AND deleted_at IS NULL;

-- name: GetUsuarioByID :one
SELECT * FROM usuarios
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateLastLogin :exec
UPDATE usuarios SET last_login_at = NOW() WHERE id = $1;

-- name: UpdatePassword :exec
UPDATE usuarios SET password_hash = $1, updated_at = NOW() WHERE id = $2;

-- name: ListUsuariosPaginated :many
SELECT * FROM usuarios
WHERE (@include_deleted::bool OR deleted_at IS NULL)
  AND (@rol::text = '' OR rol::text = @rol)
  AND (@activo::text = 'all' OR (@activo::text = 'true' AND is_active = true) OR (@activo::text = 'false' AND is_active = false))
  AND (@q::text = '' OR nombre_completo ILIKE '%' || @q || '%' OR correo ILIKE '%' || @q || '%')
ORDER BY
  CASE WHEN @sort_col::text = 'nombre_completo' AND @sort_dir::text = 'asc' THEN nombre_completo END ASC,
  CASE WHEN @sort_col::text = 'nombre_completo' AND @sort_dir::text = 'desc' THEN nombre_completo END DESC,
  CASE WHEN @sort_col::text = 'created_at' AND @sort_dir::text = 'asc' THEN created_at END ASC,
  CASE WHEN @sort_col::text = 'created_at' AND @sort_dir::text = 'desc' THEN created_at END DESC,
  CASE WHEN @sort_col::text = 'last_login_at' AND @sort_dir::text = 'asc' THEN last_login_at END ASC NULLS LAST,
  CASE WHEN @sort_col::text = 'last_login_at' AND @sort_dir::text = 'desc' THEN last_login_at END DESC NULLS LAST,
  id ASC
LIMIT $1 OFFSET $2;

-- name: CountUsuarios :one
SELECT COUNT(*) FROM usuarios
WHERE (@include_deleted::bool OR deleted_at IS NULL)
  AND (@rol::text = '' OR rol::text = @rol)
  AND (@activo::text = 'all' OR (@activo::text = 'true' AND is_active = true) OR (@activo::text = 'false' AND is_active = false))
  AND (@q::text = '' OR nombre_completo ILIKE '%' || @q || '%' OR correo ILIKE '%' || @q || '%');

-- name: CountActiveAdmins :one
SELECT COUNT(*) FROM usuarios WHERE rol = 'admin' AND is_active = true AND deleted_at IS NULL;

-- name: InsertUsuario :one
INSERT INTO usuarios (sede_id, nombre_completo, correo, password_hash, rol)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateUsuario :one
UPDATE usuarios SET
  nombre_completo = COALESCE(sqlc.narg('nombre_completo'), nombre_completo),
  correo = COALESCE(sqlc.narg('correo'), correo),
  rol = COALESCE(sqlc.narg('rol'), rol),
  updated_at = NOW()
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;

-- name: ActivateUsuario :exec
UPDATE usuarios SET is_active = true, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: DeactivateUsuario :exec
UPDATE usuarios SET is_active = false, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteUsuario :exec
UPDATE usuarios SET is_active = false, deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: ExistsCorreo :one
SELECT EXISTS(
  SELECT 1 FROM usuarios WHERE LOWER(correo) = LOWER($1)
);
