package reportes_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	reportesrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/reportes"
	usecasereportes "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/reportes"
)

type testEnv struct {
	pool        *pgxpool.Pool
	service     *usecasereportes.Service
	loc         *time.Location
	sedeID      int64
	adminID     int64
	vendedoraID int64
	efectivoID  int64
	nequiID     int64
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	pool := requireTestPool(t)

	if _, err := pool.Exec(context.Background(),
		`TRUNCATE ventas, venta_items, movimientos_inventario, stock_actual,
		 cuadres_caja, pagos_caja, consignaciones, sesiones_laborales,
		 fragancias, variantes_envase, modelos_envase, productos, metodos_pago,
		 usuarios, sedes, auditoria RESTART IDENTITY CASCADE`,
	); err != nil {
		t.Fatalf("truncating tables between tests: %v", err)
	}

	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		loc = time.FixedZone("America/Bogota", -5*60*60)
	}

	e := &testEnv{pool: pool, loc: loc}
	e.sedeID = seedSede(t, pool, "Sede Test")
	e.adminID = seedUsuario(t, pool, e.sedeID, "admin@test.local", "admin")
	e.vendedoraID = seedUsuario(t, pool, e.sedeID, "vendedora@test.local", "vendedora")
	e.efectivoID = seedMetodoPago(t, pool, "Efectivo", "efectivo")
	e.nequiID = seedMetodoPago(t, pool, "Nequi", "nequi")

	e.service = usecasereportes.NewService(reportesrepo.NewPostgres(pool), loc)
	return e
}

func (e *testEnv) exec(t *testing.T, sql string, args ...any) {
	t.Helper()
	if _, err := e.pool.Exec(context.Background(), sql, args...); err != nil {
		t.Fatalf("exec %q: %v", sql, err)
	}
}

func (e *testEnv) queryID(t *testing.T, sql string, args ...any) int64 {
	t.Helper()
	var id int64
	if err := e.pool.QueryRow(context.Background(), sql, args...).Scan(&id); err != nil {
		t.Fatalf("query %q: %v", sql, err)
	}
	return id
}

func seedSede(t *testing.T, pool *pgxpool.Pool, nombre string) int64 {
	t.Helper()
	var id int64
	if err := pool.QueryRow(context.Background(),
		`INSERT INTO sedes (nombre) VALUES ($1) RETURNING id`, nombre).Scan(&id); err != nil {
		t.Fatalf("seeding sede: %v", err)
	}
	return id
}

func seedUsuario(t *testing.T, pool *pgxpool.Pool, sedeID int64, correo, rol string) int64 {
	t.Helper()
	var id int64
	if err := pool.QueryRow(context.Background(),
		`INSERT INTO usuarios (sede_id, nombre_completo, correo, password_hash, rol)
		 VALUES ($1, 'Usuario Prueba', $2, 'x', $3) RETURNING id`,
		sedeID, correo, rol).Scan(&id); err != nil {
		t.Fatalf("seeding usuario: %v", err)
	}
	return id
}

func seedMetodoPago(t *testing.T, pool *pgxpool.Pool, nombre, codigo string) int64 {
	t.Helper()
	var id int64
	if err := pool.QueryRow(context.Background(),
		`INSERT INTO metodos_pago (nombre, codigo) VALUES ($1, $2) RETURNING id`,
		nombre, codigo).Scan(&id); err != nil {
		t.Fatalf("seeding metodo_pago: %v", err)
	}
	return id
}

func (e *testEnv) seedFragancia(t *testing.T, nombre string, minimo float64) int64 {
	return e.queryID(t,
		`INSERT INTO fragancias (sede_id, nombre_comercial, genero, gramos_minimo, numero_genero)
		 VALUES ($1, $2, 'femenina', $3,
		   (SELECT COALESCE(MAX(numero_genero), 0) + 1 FROM fragancias WHERE sede_id = $1 AND genero = 'femenina'))
		 RETURNING id`,
		e.sedeID, nombre, minimo)
}

func (e *testEnv) seedProducto(t *testing.T, nombre string, precio float64, minimo int) int64 {
	return e.queryID(t,
		`INSERT INTO productos (sede_id, nombre, categoria, precio, stock_minimo)
		 VALUES ($1, $2, 'otro', $3, $4) RETURNING id`,
		e.sedeID, nombre, precio, minimo)
}

func (e *testEnv) seedStock(t *testing.T, tipoItem string, itemID int64, ubicacion string, cantidad float64) {
	e.exec(t,
		`INSERT INTO stock_actual (sede_id, tipo_item, item_id, ubicacion, cantidad)
		 VALUES ($1, $2::tipo_item_enum, $3, $4::ubicacion_enum, $5)`,
		e.sedeID, tipoItem, itemID, ubicacion, cantidad)
}

func (e *testEnv) seedVenta(t *testing.T, metodoID int64, total float64) int64 {
	return e.queryID(t,
		`INSERT INTO ventas (sede_id, usuario_id, metodo_pago_id, subtotal, descuento_pct, descuento_monto, total)
		 VALUES ($1, $2, $3, $4, 0, 0, $4) RETURNING id`,
		e.sedeID, e.vendedoraID, metodoID, total)
}

func (e *testEnv) seedVentaItemProducto(t *testing.T, ventaID, productoID int64, cantidad int, precio float64) {
	e.exec(t,
		`INSERT INTO venta_items (venta_id, tipo_linea, producto_id, cantidad, precio_unitario, subtotal)
		 VALUES ($1, 'producto_otro', $2, $3, $4, $4)`,
		ventaID, productoID, cantidad, precio)
}

func (e *testEnv) seedMovimiento(t *testing.T, tipo, tipoItem string, itemID int64, ubicacion string, cantidad float64, motivo string, ventaID *int64) {
	e.exec(t,
		`INSERT INTO movimientos_inventario
		 (sede_id, usuario_id, tipo_item, item_id, tipo, ubicacion, cantidad, stock_anterior, stock_posterior, motivo, venta_id)
		 VALUES ($1, $2, $3::tipo_item_enum, $4, $5::tipo_movimiento_enum, $6::ubicacion_enum, $7, 0, $7, $8, $9)`,
		e.sedeID, e.vendedoraID, tipoItem, itemID, tipo, ubicacion, cantidad, motivo, ventaID)
}

func (e *testEnv) seedCuadre(t *testing.T, estado string, fecha time.Time, efectivo float64) {
	var cerradoPor *int64
	var cerradoAt *time.Time
	if estado == "cerrado" {
		cerradoPor = &e.adminID
		now := time.Now()
		cerradoAt = &now
	}
	e.exec(t,
		`INSERT INTO cuadres_caja
		 (sede_id, fecha, estado, fondo_base, total_efectivo, saldo_calculado, cerrado_por_usuario_id, cerrado_at)
		 VALUES ($1, $2::date, $3::estado_cuadre_enum, 100000, $4, $4, $5, $6)`,
		e.sedeID, fecha.Format("2006-01-02"), estado, efectivo, cerradoPor, cerradoAt)
}

func (e *testEnv) seedSesionCerrada(t *testing.T, usuarioID int64, entrada, salida time.Time) {
	e.exec(t,
		`INSERT INTO sesiones_laborales (sede_id, usuario_id, entrada_at, salida_at, horas_trabajadas)
		 VALUES ($1, $2, $3::timestamptz, $4::timestamptz, $4::timestamptz - $3::timestamptz)`,
		e.sedeID, usuarioID, entrada, salida)
}
