package auditoria_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	usecaseauditoria "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/auditoria"
)

func (e *testEnv) lastEventoID(t *testing.T) int64 {
	t.Helper()
	var id int64
	if err := e.pool.QueryRow(context.Background(),
		`SELECT id FROM auditoria ORDER BY id DESC LIMIT 1`).Scan(&id); err != nil {
		t.Fatalf("fetching last evento id: %v", err)
	}
	return id
}

func TestListOrdenDesc(t *testing.T) {
	e := newTestEnv(t)
	t0 := time.Now().Add(-time.Hour)
	e.seedEvento(t, &e.usuarioID, "login_success", "", "", t0)
	e.seedEvento(t, &e.usuarioID, "usuario_creado", "usuarios", "", t0.Add(time.Minute))
	e.seedEvento(t, &e.usuarioID, "logout", "", "", t0.Add(2*time.Minute))

	res, err := e.service.List(context.Background(), usecaseauditoria.ListInput{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if res.Total != 3 || len(res.Eventos) != 3 {
		t.Fatalf("got total=%d len=%d, want 3/3", res.Total, len(res.Eventos))
	}
	if res.Eventos[0].Accion != "logout" {
		t.Errorf("newest first: Eventos[0].Accion = %q, want 'logout'", res.Eventos[0].Accion)
	}
}

func TestListFiltroAccion(t *testing.T) {
	e := newTestEnv(t)
	t0 := time.Now().Add(-time.Hour)
	e.seedEvento(t, &e.usuarioID, "login_success", "", "", t0)
	e.seedEvento(t, &e.usuarioID, "logout", "", "", t0.Add(time.Minute))

	res, err := e.service.List(context.Background(), usecaseauditoria.ListInput{Accion: "login_success"})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if res.Total != 1 || len(res.Eventos) != 1 || res.Eventos[0].Accion != "login_success" {
		t.Fatalf("accion filter: total=%d eventos=%+v", res.Total, res.Eventos)
	}
}

func TestListFiltroUsuario(t *testing.T) {
	e := newTestEnv(t)
	t0 := time.Now().Add(-time.Hour)
	e.seedEvento(t, &e.usuarioID, "login_success", "", "", t0)
	e.seedEvento(t, &e.usuario2ID, "login_success", "", "", t0.Add(time.Minute))

	res, err := e.service.List(context.Background(), usecaseauditoria.ListInput{UsuarioID: e.usuario2ID})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if res.Total != 1 || len(res.Eventos) != 1 {
		t.Fatalf("usuario filter: total=%d len=%d", res.Total, len(res.Eventos))
	}
	if res.Eventos[0].UsuarioID == nil || *res.Eventos[0].UsuarioID != e.usuario2ID {
		t.Errorf("expected usuario_id %d, got %v", e.usuario2ID, res.Eventos[0].UsuarioID)
	}
}

func TestListFiltroTabla(t *testing.T) {
	e := newTestEnv(t)
	t0 := time.Now().Add(-time.Hour)
	e.seedEvento(t, &e.usuarioID, "usuario_creado", "usuarios", "", t0)
	e.seedEvento(t, &e.usuarioID, "fragancia_creada", "fragancias", "", t0.Add(time.Minute))

	res, err := e.service.List(context.Background(), usecaseauditoria.ListInput{TablaAfectada: "fragancias"})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if res.Total != 1 || len(res.Eventos) != 1 || res.Eventos[0].Accion != "fragancia_creada" {
		t.Fatalf("tabla filter: total=%d eventos=%+v", res.Total, res.Eventos)
	}
}

func TestListFiltroRango(t *testing.T) {
	e := newTestEnv(t)
	t0 := time.Now().Add(-time.Hour)
	e.seedEvento(t, &e.usuarioID, "a", "", "", t0)
	e.seedEvento(t, &e.usuarioID, "b", "", "", t0.Add(time.Minute))
	e.seedEvento(t, &e.usuarioID, "c", "", "", t0.Add(2*time.Minute))

	desde := t0.Add(30 * time.Second)
	hasta := t0.Add(90 * time.Second)
	res, err := e.service.List(context.Background(), usecaseauditoria.ListInput{FechaDesde: &desde, FechaHasta: &hasta})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if res.Total != 1 || len(res.Eventos) != 1 || res.Eventos[0].Accion != "b" {
		t.Fatalf("rango filter: total=%d eventos=%+v", res.Total, res.Eventos)
	}
}

func TestListPaginacion(t *testing.T) {
	e := newTestEnv(t)
	t0 := time.Now().Add(-time.Hour)
	for i := 0; i < 3; i++ {
		e.seedEvento(t, &e.usuarioID, "login_success", "", "", t0.Add(time.Duration(i)*time.Minute))
	}

	page1, err := e.service.List(context.Background(), usecaseauditoria.ListInput{Page: 1, PageSize: 2})
	if err != nil {
		t.Fatalf("List page1: %v", err)
	}
	if page1.Total != 3 || len(page1.Eventos) != 2 || page1.PageSize != 2 {
		t.Fatalf("page1: total=%d len=%d size=%d", page1.Total, len(page1.Eventos), page1.PageSize)
	}
	page2, err := e.service.List(context.Background(), usecaseauditoria.ListInput{Page: 2, PageSize: 2})
	if err != nil {
		t.Fatalf("List page2: %v", err)
	}
	if len(page2.Eventos) != 1 {
		t.Fatalf("page2: len=%d, want 1", len(page2.Eventos))
	}
}

func TestGetNotFound(t *testing.T) {
	e := newTestEnv(t)

	_, err := e.service.Get(context.Background(), 999999)
	if err == nil {
		t.Fatal("expected not-found error, got nil")
	}
	var de *domainerrors.DomainError
	if !errors.As(err, &de) || de.Code != domainerrors.CodeNotFound {
		t.Fatalf("expected CodeNotFound DomainError, got %v", err)
	}
}

func TestGetJSONParseado(t *testing.T) {
	e := newTestEnv(t)
	e.seedEvento(t, &e.usuarioID, "usuario_editado", "usuarios", `{"nombre":"Ana","activo":true}`, time.Now())

	ev, err := e.service.Get(context.Background(), e.lastEventoID(t))
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	var parsed map[string]any
	if err := json.Unmarshal(ev.DatosDespues, &parsed); err != nil {
		t.Fatalf("datos_despues is not valid JSON: %v (raw=%s)", err, ev.DatosDespues)
	}
	if parsed["nombre"] != "Ana" {
		t.Errorf("datos_despues.nombre = %v, want 'Ana'", parsed["nombre"])
	}
	if ev.IP == nil || *ev.IP != "192.168.1.10" {
		t.Errorf("ip = %v, want 192.168.1.10", ev.IP)
	}
	if ev.Usuario == nil || ev.Usuario.ID != e.usuarioID {
		t.Errorf("usuario brief = %v, want id %d", ev.Usuario, e.usuarioID)
	}
}
