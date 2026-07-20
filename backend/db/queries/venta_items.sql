-- name: InsertVentaItem :one
INSERT INTO venta_items (
    venta_id, tipo_linea,
    fragancia_id, variante_envase_id, producto_id, feromona_producto_id,
    gramos_fragancia, cantidad, precio_unitario, subtotal
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetVentaItemsByVentaID :many
SELECT
    vi.*,
    f.nombre_comercial AS fragancia_nombre,
    f.nombre_alternativo AS fragancia_alternativo,
    me.tipo AS modelo_envase_tipo,
    me.tamano_oz AS modelo_envase_tamano_oz,
    ve.color AS variante_envase_color,
    p.nombre AS producto_nombre,
    p.categoria AS producto_categoria,
    fp.nombre AS feromona_nombre
FROM venta_items vi
LEFT JOIN fragancias f ON f.id = vi.fragancia_id
LEFT JOIN variantes_envase ve ON ve.id = vi.variante_envase_id
LEFT JOIN modelos_envase me ON me.id = ve.modelo_envase_id
LEFT JOIN productos p ON p.id = vi.producto_id
LEFT JOIN productos fp ON fp.id = vi.feromona_producto_id
WHERE vi.venta_id = $1
ORDER BY vi.id ASC;
