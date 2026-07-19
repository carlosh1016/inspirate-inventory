package usuarios_test

import (
	"context"
	"testing"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/usuarios"
)

const unknownID = 9_999_999

func TestGetUnknownUserNotFound(t *testing.T) {
	env := newTestEnv(t)
	_, err := env.service.Get(context.Background(), unknownID)
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestActivateUnknownUserNotFound(t *testing.T) {
	env := newTestEnv(t)
	adminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-admin"), "admin", true)

	err := env.service.Activate(context.Background(), usecase.ActivateInput{TargetID: unknownID, RequesterID: adminID})
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestDeactivateUnknownUserNotFound(t *testing.T) {
	env := newTestEnv(t)
	adminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-admin"), "admin", true)

	err := env.service.Deactivate(context.Background(), usecase.DeactivateInput{TargetID: unknownID, RequesterID: adminID})
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestDeleteUnknownUserNotFound(t *testing.T) {
	env := newTestEnv(t)
	adminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-admin"), "admin", true)

	err := env.service.Delete(context.Background(), usecase.DeleteInput{TargetID: unknownID, RequesterID: adminID})
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestUpdateUnknownUserNotFound(t *testing.T) {
	env := newTestEnv(t)
	adminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-admin"), "admin", true)

	_, err := env.service.Update(context.Background(), usecase.UpdateInput{
		TargetID: unknownID, NombreCompleto: strPtr("Nuevo Nombre"), RequesterID: adminID,
	})
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestUpdatePasswordUnknownUserNotFound(t *testing.T) {
	env := newTestEnv(t)
	adminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-admin"), "admin", true)

	err := env.service.UpdatePassword(context.Background(), usecase.UpdatePasswordInput{
		TargetID: unknownID, PasswordNueva: "nueva-password-123",
		RequesterID: adminID, RequesterIsAdmin: true,
	})
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestUpdateNombreCompletoOnly(t *testing.T) {
	env := newTestEnv(t)
	adminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-admin"), "admin", true)
	targetID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-target"), "vendedora", true)

	updated, err := env.service.Update(context.Background(), usecase.UpdateInput{
		TargetID: targetID, NombreCompleto: strPtr("Nombre Actualizado"), RequesterID: adminID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.NombreCompleto != "Nombre Actualizado" {
		t.Errorf("expected updated nombre_completo, got %q", updated.NombreCompleto)
	}
}

func TestUpdateCorreoConflictWithExisting(t *testing.T) {
	env := newTestEnv(t)
	adminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-admin"), "admin", true)
	takenCorreo := uniqueCorreo(t, "-taken")
	seedUsuario(t, env.pool, env.sedeID, takenCorreo, "vendedora", true)
	targetID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-target"), "vendedora", true)

	_, err := env.service.Update(context.Background(), usecase.UpdateInput{
		TargetID: targetID, Correo: &takenCorreo, RequesterID: adminID,
	})
	assertCode(t, err, domainerrors.CodeConflict)
}

func TestListSortFallsBackOnInvalidColumn(t *testing.T) {
	env := newTestEnv(t)
	seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-a"), "vendedora", true)

	if _, err := env.service.List(context.Background(), usecase.ListInput{Sort: "correo:asc"}); err != nil {
		t.Fatalf("unexpected error with a disallowed sort column: %v", err)
	}
}

func TestListSortFallsBackOnInvalidDirection(t *testing.T) {
	env := newTestEnv(t)
	seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-a"), "vendedora", true)

	if _, err := env.service.List(context.Background(), usecase.ListInput{Sort: "nombre_completo:sideways"}); err != nil {
		t.Fatalf("unexpected error with an invalid sort direction: %v", err)
	}
}

func TestListSortByNombreCompletoAscending(t *testing.T) {
	env := newTestEnv(t)
	seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-a"), "vendedora", true)

	result, err := env.service.List(context.Background(), usecase.ListInput{Sort: "nombre_completo:asc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 1 {
		t.Errorf("expected 1 result, got %d", result.Total)
	}
}

func TestListDefaultsActivoToTrue(t *testing.T) {
	env := newTestEnv(t)
	seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-active"), "vendedora", true)
	seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-inactive"), "vendedora", false)

	result, err := env.service.List(context.Background(), usecase.ListInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 1 {
		t.Errorf("expected only the active user by default, got total=%d", result.Total)
	}
}
