-- name: InsertStockActual :exec
INSERT INTO stock_actual (sede_id, tipo_item, item_id, ubicacion, cantidad)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (sede_id, tipo_item, item_id, ubicacion) DO NOTHING;

-- name: GetStockTotalByItem :one
SELECT
  COALESCE(SUM(CASE WHEN ubicacion = 'vitrina' THEN cantidad ELSE 0 END), 0)::numeric AS vitrina,
  COALESCE(SUM(CASE WHEN ubicacion = 'bodega' THEN cantidad ELSE 0 END), 0)::numeric AS bodega
FROM stock_actual
WHERE sede_id = $1 AND tipo_item = $2 AND item_id = $3;
