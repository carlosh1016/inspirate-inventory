-- name: InsertConteo :one
INSERT INTO conteos_inventario (sede_id, periodo, creado_por_usuario_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GenerarItemsConteo :exec
-- Crea una fila por cada fragancia activa de la sede. gramos_sistema sale del
-- stock_actual agregado (bodega+vitrina). saldo_inicial es la base del
-- periodo (el gramos_fisico del conteo cerrado del mes anterior, o —si es el
-- primer conteo de esa fragancia— el stock reconstruido al inicio del mes
-- restando del stock actual todos los movimientos de este mes, ya que todo
-- movimiento en movimientos_inventario guarda un delta con signo) MÁS las
-- entradas de mercancía (compras) registradas durante este mismo mes: si se
-- compra más producto a mitad de mes, esa cantidad se suma al saldo inicial
-- en vez de perderse en el stock actual, para que el % vendido siga
-- reflejando cuánto se vendió sobre el total realmente disponible ese mes.
INSERT INTO conteo_inventario_items (conteo_id, fragancia_id, saldo_inicial, gramos_sistema)
SELECT
  @conteo_id::bigint,
  f.id,
  COALESCE(
    prev.gramos_fisico,
    GREATEST(COALESCE(stock.total, 0) - COALESCE(mov.neto, 0), 0)
  ) + COALESCE(mov.entradas, 0),
  COALESCE(stock.total, 0)
FROM fragancias f
LEFT JOIN LATERAL (
  SELECT SUM(sa.cantidad) AS total
  FROM stock_actual sa
  WHERE sa.tipo_item = 'fragancia' AND sa.item_id = f.id AND sa.sede_id = @sede_id::bigint
) stock ON true
LEFT JOIN LATERAL (
  SELECT cii.gramos_fisico
  FROM conteo_inventario_items cii
  INNER JOIN conteos_inventario ci ON ci.id = cii.conteo_id
  WHERE ci.sede_id = @sede_id::bigint
    AND ci.estado = 'cerrado'
    AND ci.periodo < @periodo::date
    AND cii.fragancia_id = f.id
    AND cii.gramos_fisico IS NOT NULL
  ORDER BY ci.periodo DESC
  LIMIT 1
) prev ON true
LEFT JOIN LATERAL (
  SELECT
    SUM(mi.cantidad) AS neto,
    SUM(mi.cantidad) FILTER (WHERE mi.tipo = 'entrada_mercancia') AS entradas
  FROM movimientos_inventario mi
  WHERE mi.tipo_item = 'fragancia' AND mi.item_id = f.id AND mi.sede_id = @sede_id::bigint
    AND mi.created_at >= @periodo::date
) mov ON true
WHERE f.sede_id = @sede_id::bigint AND f.deleted_at IS NULL AND f.activo = true;

-- name: GetConteoBySedePeriodo :one
SELECT * FROM conteos_inventario WHERE sede_id = $1 AND periodo = $2;

-- name: GetConteoByID :one
SELECT
  c.*,
  creado.nombre_completo AS creado_por_nombre,
  cerrado.nombre_completo AS cerrado_por_nombre
FROM conteos_inventario c
INNER JOIN usuarios creado ON creado.id = c.creado_por_usuario_id
LEFT JOIN usuarios cerrado ON cerrado.id = c.cerrado_por_usuario_id
WHERE c.id = $1;

-- name: ListConteoItems :many
SELECT
  cii.*,
  f.nombre_comercial AS fragancia_nombre
FROM conteo_inventario_items cii
INNER JOIN fragancias f ON f.id = cii.fragancia_id
WHERE cii.conteo_id = $1
ORDER BY f.nombre_comercial ASC;

-- name: UpdateConteoItemFisico :one
UPDATE conteo_inventario_items cii
SET gramos_fisico = $3, updated_at = NOW()
WHERE cii.id = $1 AND cii.conteo_id = $2
RETURNING cii.*, (SELECT f.nombre_comercial FROM fragancias f WHERE f.id = cii.fragancia_id) AS fragancia_nombre;

-- name: CerrarConteo :one
UPDATE conteos_inventario
SET estado = 'cerrado', cerrado_por_usuario_id = $2, cerrado_at = NOW(), updated_at = NOW()
WHERE id = $1 AND estado = 'abierto'
RETURNING *;

-- name: ListConteosPaginated :many
SELECT
  c.*,
  creado.nombre_completo AS creado_por_nombre,
  cerrado.nombre_completo AS cerrado_por_nombre
FROM conteos_inventario c
INNER JOIN usuarios creado ON creado.id = c.creado_por_usuario_id
LEFT JOIN usuarios cerrado ON cerrado.id = c.cerrado_por_usuario_id
WHERE c.sede_id = @sede_id::bigint
ORDER BY c.periodo DESC
LIMIT $1 OFFSET $2;

-- name: CountConteos :one
SELECT COUNT(*) FROM conteos_inventario WHERE sede_id = @sede_id::bigint;
