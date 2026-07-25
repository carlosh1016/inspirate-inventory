-- Queries dedicadas de agregación pesada para los reportes XLSX (M13).
-- Todas escopadas por sede_id. Las expresiones computadas/COALESCE llevan
-- cast explícito ::tipo para no romper la inferencia de sqlc (lección
-- recurrente: CASE/CONCAT/COALESCE sin cast generan interface{} o tipos no
-- nulos que panican al escanear NULL).

-- ===========================================================================
-- VENTAS
-- ===========================================================================

-- name: ReporteVentasResumen :one
-- Un solo row con los totales globales + el desglose por método de pago,
-- mismo bucketing que GetResumenVentasHoy pero sobre un rango arbitrario.
WITH ventas_rango AS (
    SELECT * FROM ventas v_raw
    WHERE v_raw.sede_id = @sede_id::bigint
      AND v_raw.created_at >= @fecha_desde::timestamptz
      AND v_raw.created_at <= @fecha_hasta::timestamptz
      AND (@usuario_id::bigint = 0 OR v_raw.usuario_id = @usuario_id)
),
totales_por_metodo AS (
    SELECT mp.codigo, COALESCE(SUM(v.total), 0)::numeric AS total_metodo
    FROM metodos_pago mp
    LEFT JOIN ventas_rango v ON v.metodo_pago_id = mp.id
    GROUP BY mp.codigo
)
SELECT
    COALESCE((SELECT SUM(total_metodo) FROM totales_por_metodo WHERE codigo = 'efectivo'), 0)::numeric AS total_efectivo,
    COALESCE((SELECT SUM(total_metodo) FROM totales_por_metodo WHERE codigo = 'nequi'), 0)::numeric AS total_nequi,
    COALESCE((SELECT SUM(total_metodo) FROM totales_por_metodo WHERE codigo = 'daviplata'), 0)::numeric AS total_daviplata,
    COALESCE((SELECT SUM(total_metodo) FROM totales_por_metodo WHERE codigo = 'transferencia'), 0)::numeric AS total_transferencia,
    COALESCE((SELECT SUM(total_metodo) FROM totales_por_metodo WHERE codigo NOT IN ('efectivo','nequi','daviplata','transferencia')), 0)::numeric AS total_otros,
    (SELECT COUNT(*) FROM ventas_rango)::bigint AS ventas_count,
    COALESCE((SELECT SUM(total) FROM ventas_rango), 0)::numeric AS total_ventas,
    COALESCE((SELECT SUM(descuento_monto) FROM ventas_rango), 0)::numeric AS descuento_total;

-- name: ReporteVentasPorVendedora :many
SELECT u.id AS usuario_id, u.nombre_completo,
       COUNT(v.id)::bigint AS ventas_count,
       COALESCE(SUM(v.total), 0)::numeric AS total
FROM ventas v
INNER JOIN usuarios u ON u.id = v.usuario_id
WHERE v.sede_id = @sede_id::bigint
  AND v.created_at >= @fecha_desde::timestamptz
  AND v.created_at <= @fecha_hasta::timestamptz
  AND (@usuario_id::bigint = 0 OR v.usuario_id = @usuario_id)
GROUP BY u.id, u.nombre_completo
ORDER BY total DESC;

-- name: ReporteVentasDetalle :many
SELECT v.id, v.created_at,
       u.nombre_completo AS usuario_nombre,
       mp.nombre AS metodo_pago_nombre,
       v.subtotal, v.descuento_pct, v.descuento_monto, v.total, v.observaciones
FROM ventas v
INNER JOIN usuarios u ON u.id = v.usuario_id
INNER JOIN metodos_pago mp ON mp.id = v.metodo_pago_id
WHERE v.sede_id = @sede_id::bigint
  AND v.created_at >= @fecha_desde::timestamptz
  AND v.created_at <= @fecha_hasta::timestamptz
  AND (@usuario_id::bigint = 0 OR v.usuario_id = @usuario_id)
ORDER BY v.created_at ASC, v.id ASC;

-- name: ReporteVentasItems :many
-- Detalle línea a línea. fragancia/producto/feromona son columnas directas de
-- tablas LEFT-JOIN (sqlc las infiere nullable → pgtype.Text). envase_nombre es
-- un CONCAT: se envuelve en COALESCE(...,'')::text para evitar el panic de
-- "scan NULL into string" (sale '' cuando la línea no tiene envase).
SELECT
    vi.venta_id,
    v.created_at,
    vi.tipo_linea::text AS tipo_linea,
    f.nombre_comercial AS fragancia_nombre,
    COALESCE(
        CASE WHEN vi.variante_envase_id IS NOT NULL
             THEN CONCAT(me.tipo, ' ', me.tamano_oz::text, 'oz ', ve.color)
        END, '')::text AS envase_nombre,
    p.nombre AS producto_nombre,
    fp.nombre AS feromona_nombre,
    vi.gramos_fragancia,
    vi.cantidad,
    vi.precio_unitario,
    vi.subtotal
FROM venta_items vi
INNER JOIN ventas v ON v.id = vi.venta_id
LEFT JOIN fragancias f ON f.id = vi.fragancia_id
LEFT JOIN variantes_envase ve ON ve.id = vi.variante_envase_id
LEFT JOIN modelos_envase me ON me.id = ve.modelo_envase_id
LEFT JOIN productos p ON p.id = vi.producto_id
LEFT JOIN productos fp ON fp.id = vi.feromona_producto_id
WHERE v.sede_id = @sede_id::bigint
  AND v.created_at >= @fecha_desde::timestamptz
  AND v.created_at <= @fecha_hasta::timestamptz
  AND (@usuario_id::bigint = 0 OR v.usuario_id = @usuario_id)
ORDER BY vi.venta_id ASC, vi.id ASC;

-- ===========================================================================
-- STOCK (snapshot actual, sin rango)
-- ===========================================================================

-- name: ReporteStockFragancias :many
WITH agg AS (
    SELECT item_id,
           SUM(CASE WHEN ubicacion = 'vitrina' THEN cantidad ELSE 0 END) AS vit,
           SUM(CASE WHEN ubicacion = 'bodega' THEN cantidad ELSE 0 END) AS bod,
           SUM(cantidad) AS tot
    FROM stock_actual
    WHERE sede_id = @sede_id::bigint AND tipo_item = 'fragancia'
    GROUP BY item_id
)
SELECT f.nombre_comercial, f.nombre_alternativo, f.genero::text AS genero,
       COALESCE(a.vit, 0)::numeric AS stock_vitrina,
       COALESCE(a.bod, 0)::numeric AS stock_bodega,
       COALESCE(a.tot, 0)::numeric AS stock_total,
       f.gramos_minimo AS minimo,
       (COALESCE(a.tot, 0) < f.gramos_minimo) AS bajo_minimo
FROM fragancias f
LEFT JOIN agg a ON a.item_id = f.id
WHERE f.deleted_at IS NULL AND (@include_inactivos::bool OR f.activo = true)
ORDER BY f.nombre_comercial ASC;

-- name: ReporteStockEnvases :many
WITH agg AS (
    SELECT item_id,
           SUM(CASE WHEN ubicacion = 'vitrina' THEN cantidad ELSE 0 END) AS vit,
           SUM(CASE WHEN ubicacion = 'bodega' THEN cantidad ELSE 0 END) AS bod,
           SUM(cantidad) AS tot
    FROM stock_actual
    WHERE sede_id = @sede_id::bigint AND tipo_item = 'variante_envase'
    GROUP BY item_id
)
SELECT me.tipo, me.tamano_oz, ve.color,
       me.precio_solo, me.precio_con_fragancia, me.precio_recarga,
       COALESCE(a.vit, 0)::numeric AS stock_vitrina,
       COALESCE(a.bod, 0)::numeric AS stock_bodega,
       COALESCE(a.tot, 0)::numeric AS stock_total,
       ve.stock_minimo::numeric AS minimo,
       (COALESCE(a.tot, 0) < ve.stock_minimo) AS bajo_minimo
FROM variantes_envase ve
INNER JOIN modelos_envase me ON me.id = ve.modelo_envase_id
LEFT JOIN agg a ON a.item_id = ve.id
WHERE ve.deleted_at IS NULL AND (@include_inactivos::bool OR ve.activo = true)
ORDER BY me.tipo ASC, ve.color ASC;

-- name: ReporteStockProductos :many
WITH agg AS (
    SELECT item_id,
           SUM(CASE WHEN ubicacion = 'vitrina' THEN cantidad ELSE 0 END) AS vit,
           SUM(CASE WHEN ubicacion = 'bodega' THEN cantidad ELSE 0 END) AS bod,
           SUM(cantidad) AS tot
    FROM stock_actual
    WHERE sede_id = @sede_id::bigint AND tipo_item = 'producto'
    GROUP BY item_id
)
SELECT p.nombre, p.categoria::text AS categoria, p.precio,
       COALESCE(a.vit, 0)::numeric AS stock_vitrina,
       COALESCE(a.bod, 0)::numeric AS stock_bodega,
       COALESCE(a.tot, 0)::numeric AS stock_total,
       p.stock_minimo::numeric AS minimo,
       (COALESCE(a.tot, 0) < p.stock_minimo) AS bajo_minimo
FROM productos p
LEFT JOIN agg a ON a.item_id = p.id
WHERE p.deleted_at IS NULL AND (@include_inactivos::bool OR p.activo = true)
ORDER BY p.nombre ASC;

-- name: ReporteStockAlertas :many
-- Solo items bajo mínimo (total < mínimo). El faltante es mínimo - total. La
-- columna "Ubicación" del reporte se rellena en Go con "Total" (el mínimo es
-- un umbral sobre el stock total, no por ubicación).
WITH agg AS (
    SELECT tipo_item, item_id, SUM(cantidad) AS tot
    FROM stock_actual
    WHERE sede_id = @sede_id::bigint
    GROUP BY tipo_item, item_id
)
SELECT * FROM (
    SELECT 'Fragancia'::text AS tipo, f.nombre_comercial AS nombre,
           COALESCE(a.tot, 0)::numeric AS stock_actual,
           f.gramos_minimo AS minimo,
           (f.gramos_minimo - COALESCE(a.tot, 0))::numeric AS faltante
    FROM fragancias f
    LEFT JOIN agg a ON a.tipo_item = 'fragancia' AND a.item_id = f.id
    WHERE f.deleted_at IS NULL AND (@include_inactivos::bool OR f.activo = true)
      AND COALESCE(a.tot, 0) < f.gramos_minimo

    UNION ALL

    SELECT 'Envase'::text, CONCAT(me.tipo, ' ', me.tamano_oz::text, 'oz ', ve.color),
           COALESCE(a.tot, 0)::numeric, ve.stock_minimo::numeric,
           (ve.stock_minimo - COALESCE(a.tot, 0))::numeric
    FROM variantes_envase ve
    INNER JOIN modelos_envase me ON me.id = ve.modelo_envase_id
    LEFT JOIN agg a ON a.tipo_item = 'variante_envase' AND a.item_id = ve.id
    WHERE ve.deleted_at IS NULL AND (@include_inactivos::bool OR ve.activo = true)
      AND COALESCE(a.tot, 0) < ve.stock_minimo

    UNION ALL

    SELECT 'Producto'::text, p.nombre,
           COALESCE(a.tot, 0)::numeric, p.stock_minimo::numeric,
           (p.stock_minimo - COALESCE(a.tot, 0))::numeric
    FROM productos p
    LEFT JOIN agg a ON a.tipo_item = 'producto' AND a.item_id = p.id
    WHERE p.deleted_at IS NULL AND (@include_inactivos::bool OR p.activo = true)
      AND COALESCE(a.tot, 0) < p.stock_minimo
) alertas
ORDER BY tipo ASC, nombre ASC;

-- ===========================================================================
-- MOVIMIENTOS
-- ===========================================================================

-- name: ReporteMovimientos :many
-- Igual a ListMovimientosPaginated pero sin LIMIT/OFFSET (el reporte descarga
-- todo el rango). Mismos filtros y mismo item_nombre calculado.
SELECT m.*,
  (CASE m.tipo_item
    WHEN 'fragancia' THEN (SELECT nombre_comercial FROM fragancias WHERE id = m.item_id)
    WHEN 'variante_envase' THEN (
      SELECT CONCAT(me.tipo, ' ', me.tamano_oz::text, 'oz ', ve.color)
      FROM variantes_envase ve
      JOIN modelos_envase me ON me.id = ve.modelo_envase_id
      WHERE ve.id = m.item_id
    )
    WHEN 'producto' THEN (SELECT nombre FROM productos WHERE id = m.item_id)
  END)::text AS item_nombre,
  u.nombre_completo AS usuario_nombre
FROM movimientos_inventario m
INNER JOIN usuarios u ON u.id = m.usuario_id
WHERE
  m.sede_id = @sede_id::bigint
  AND (@tipo_item::text = '' OR m.tipo_item::text = @tipo_item)
  AND (@item_id::bigint = 0 OR m.item_id = @item_id)
  AND (@tipo::text = '' OR m.tipo::text = @tipo)
  AND (@usuario_id::bigint = 0 OR m.usuario_id = @usuario_id)
  AND (@fecha_desde::timestamptz IS NULL OR m.created_at >= @fecha_desde)
  AND (@fecha_hasta::timestamptz IS NULL OR m.created_at <= @fecha_hasta)
ORDER BY m.created_at DESC, m.id DESC;

-- ===========================================================================
-- CUADRES DE CAJA (solo cerrados)
-- ===========================================================================

-- name: ReporteCuadresCerrados :many
SELECT c.fecha, c.estado::text AS estado, c.fondo_base,
       c.total_efectivo, c.total_nequi, c.total_daviplata, c.total_transferencia, c.total_otros,
       c.total_pagos, c.total_consignaciones, c.valor_turno, c.saldo_calculado,
       u.nombre_completo AS cerrado_por, c.cerrado_at, c.observaciones
FROM cuadres_caja c
LEFT JOIN usuarios u ON u.id = c.cerrado_por_usuario_id
WHERE c.sede_id = @sede_id::bigint
  AND c.estado = 'cerrado'
  AND c.fecha >= @fecha_desde::date
  AND c.fecha <= @fecha_hasta::date
ORDER BY c.fecha ASC;

-- name: ReporteCuadresPagos :many
SELECT c.fecha AS fecha_cuadre, pc.concepto, pc.monto,
       u.nombre_completo AS usuario_nombre, pc.created_at
FROM pagos_caja pc
INNER JOIN cuadres_caja c ON c.id = pc.cuadre_caja_id
INNER JOIN usuarios u ON u.id = pc.usuario_id
WHERE c.sede_id = @sede_id::bigint
  AND c.estado = 'cerrado'
  AND c.fecha >= @fecha_desde::date
  AND c.fecha <= @fecha_hasta::date
ORDER BY c.fecha ASC, pc.created_at ASC;

-- name: ReporteCuadresConsignaciones :many
SELECT c.fecha AS fecha_cuadre, cons.monto, cons.banco, cons.referencia,
       u.nombre_completo AS usuario_nombre, cons.created_at
FROM consignaciones cons
INNER JOIN cuadres_caja c ON c.id = cons.cuadre_caja_id
INNER JOIN usuarios u ON u.id = cons.usuario_id
WHERE c.sede_id = @sede_id::bigint
  AND c.estado = 'cerrado'
  AND c.fecha >= @fecha_desde::date
  AND c.fecha <= @fecha_hasta::date
ORDER BY c.fecha ASC, cons.created_at ASC;

-- ===========================================================================
-- SESIONES LABORALES
-- ===========================================================================

-- name: ReporteSesionesResumen :many
-- COALESCE(SUM(interval), INTERVAL '0')::interval: el cast ::interval es
-- obligatorio o sqlc genera interface{} (lección de Tanda 5). dias_trabajados
-- cuenta fechas de entrada distintas en zona Colombia.
SELECT u.id AS usuario_id, u.nombre_completo,
       COALESCE(SUM(s.horas_trabajadas), INTERVAL '0')::interval AS total_horas,
       COUNT(DISTINCT (s.entrada_at AT TIME ZONE 'America/Bogota')::date)::bigint AS dias_trabajados,
       COUNT(*)::bigint AS sesiones_count
FROM sesiones_laborales s
INNER JOIN usuarios u ON u.id = s.usuario_id
WHERE s.sede_id = @sede_id::bigint
  AND s.entrada_at >= @fecha_desde::timestamptz
  AND s.entrada_at <= @fecha_hasta::timestamptz
  AND (@usuario_id::bigint = 0 OR s.usuario_id = @usuario_id)
GROUP BY u.id, u.nombre_completo
ORDER BY u.nombre_completo ASC;

-- name: ReporteSesionesDetalle :many
SELECT u.nombre_completo, s.entrada_at, s.salida_at, s.horas_trabajadas
FROM sesiones_laborales s
INNER JOIN usuarios u ON u.id = s.usuario_id
WHERE s.sede_id = @sede_id::bigint
  AND s.entrada_at >= @fecha_desde::timestamptz
  AND s.entrada_at <= @fecha_hasta::timestamptz
  AND (@usuario_id::bigint = 0 OR s.usuario_id = @usuario_id)
ORDER BY u.nombre_completo ASC, s.entrada_at ASC;
