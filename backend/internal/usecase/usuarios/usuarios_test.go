package usuarios_test

import (
	"context"
	"errors"
	"testing"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	refreshtokens "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/refresh_tokens"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/usuarios"
)

func strPtr(s string) *string { return &s }

func TestCreateSuccess(t *testing.T) {
	env := newTestEnv(t)
	adminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-admin"), "admin", true)
	correo := uniqueCorreo(t, "-nueva")

	user, err := env.service.Create(context.Background(), usecase.CreateInput{
		SedeID: env.sedeID, NombreCompleto: "Nueva Vendedora", Correo: correo,
		Password: "password123", Rol: "vendedora", RequesterID: adminID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Correo != correo || string(user.Rol) != "vendedora" {
		t.Errorf("unexpected user: %+v", user)
	}
}

func TestCreateDuplicateCorreoConflicts(t *testing.T) {
	env := newTestEnv(t)
	adminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-admin"), "admin", true)
	correo := uniqueCorreo(t, "-dup")
	seedUsuario(t, env.pool, env.sedeID, correo, "vendedora", true)

	_, err := env.service.Create(context.Background(), usecase.CreateInput{
		SedeID: env.sedeID, NombreCompleto: "Otra Persona", Correo: correo,
		Password: "password123", Rol: "vendedora", RequesterID: adminID,
	})
	assertCode(t, err, domainerrors.CodeConflict)
}

func TestUpdateOwnRoleRejected(t *testing.T) {
	env := newTestEnv(t)
	adminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-admin"), "admin", true)

	_, err := env.service.Update(context.Background(), usecase.UpdateInput{
		TargetID: adminID, Rol: strPtr("vendedora"), RequesterID: adminID,
	})
	assertCode(t, err, domainerrors.CodeValidation)
}

func TestUpdateOtherAdminRoleAllowedWithQuorum(t *testing.T) {
	env := newTestEnv(t)
	adminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-admin1"), "admin", true)
	otherAdminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-admin2"), "admin", true)

	updated, err := env.service.Update(context.Background(), usecase.UpdateInput{
		TargetID: otherAdminID, Rol: strPtr("vendedora"), RequesterID: adminID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(updated.Rol) != "vendedora" {
		t.Errorf("expected rol=vendedora, got %q", updated.Rol)
	}
}

func TestUpdateDemotingLastAdminRejected(t *testing.T) {
	env := newTestEnv(t)
	adminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-admin1"), "admin", true)
	onlyOtherAdminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-admin2"), "admin", true)

	// Demote the requester's peer first isn't needed; demote onlyOtherAdminID
	// down to 1 active admin (adminID) then try to demote adminID too via a
	// second requester... simplest: deactivate onlyOtherAdminID so only
	// adminID remains, then try demoting adminID via a third admin.
	if err := env.service.Deactivate(context.Background(), usecase.DeactivateInput{
		TargetID: onlyOtherAdminID, RequesterID: adminID,
	}); err != nil {
		t.Fatalf("setup: deactivate failed: %v", err)
	}

	thirdAdminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-admin3"), "admin", true)

	_, err := env.service.Update(context.Background(), usecase.UpdateInput{
		TargetID: adminID, Rol: strPtr("vendedora"), RequesterID: thirdAdminID,
	})
	if err != nil {
		t.Fatalf("unexpected error demoting one of two remaining admins: %v", err)
	}

	// Now only thirdAdminID is an active admin; demoting them must fail.
	_, err = env.service.Update(context.Background(), usecase.UpdateInput{
		TargetID: thirdAdminID, Rol: strPtr("vendedora"), RequesterID: adminID,
	})
	assertCode(t, err, domainerrors.CodeBusinessRule)
}

func TestUpdateCorreoRevokesSessions(t *testing.T) {
	env := newTestEnv(t)
	adminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-admin"), "admin", true)
	targetID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-target"), "vendedora", true)

	if _, err := env.pool.Exec(context.Background(),
		`INSERT INTO refresh_tokens (usuario_id, token_hash, expires_at) VALUES ($1, $2, NOW() + INTERVAL '1 day')`,
		targetID, sha256Hex(uniqueCorreo(t, "-token")),
	); err != nil {
		t.Fatalf("seeding refresh token: %v", err)
	}

	newCorreo := uniqueCorreo(t, "-changed")
	if _, err := env.service.Update(context.Background(), usecase.UpdateInput{
		TargetID: targetID, Correo: &newCorreo, RequesterID: adminID,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rtRepo := refreshtokens.NewPostgres(env.pool)
	rt, err := rtRepo.GetByHash(context.Background(), sha256Hex(uniqueCorreo(t, "-token")))
	if err != nil {
		t.Fatalf("fetching refresh token: %v", err)
	}
	if rt.RevokedAt.Time.IsZero() || !rt.RevokedAt.Valid {
		t.Error("expected the refresh token to be revoked after a correo change")
	}
}

func TestDeleteSelfRejected(t *testing.T) {
	env := newTestEnv(t)
	adminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-admin"), "admin", true)

	err := env.service.Delete(context.Background(), usecase.DeleteInput{TargetID: adminID, RequesterID: adminID})
	assertCode(t, err, domainerrors.CodeBusinessRule)
}

func TestDeactivateSelfRejected(t *testing.T) {
	env := newTestEnv(t)
	adminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-admin"), "admin", true)

	err := env.service.Deactivate(context.Background(), usecase.DeactivateInput{TargetID: adminID, RequesterID: adminID})
	assertCode(t, err, domainerrors.CodeBusinessRule)
}

func TestDeleteLastAdminRejected(t *testing.T) {
	env := newTestEnv(t)
	requesterID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-requester"), "admin", true)
	lastAdminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-lastadmin"), "admin", true)

	// Deactivate requesterID first so lastAdminID is the sole active admin.
	if err := env.service.Deactivate(context.Background(), usecase.DeactivateInput{
		TargetID: requesterID, RequesterID: lastAdminID,
	}); err != nil {
		t.Fatalf("setup: deactivate failed: %v", err)
	}

	// Someone else (not an admin — the rule only cares about the target's
	// quorum, not the requester's role) tries to delete the now-only active
	// admin. Using a non-admin keeps the active-admin count at 1.
	otherID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-other"), "vendedora", true)
	err := env.service.Delete(context.Background(), usecase.DeleteInput{TargetID: lastAdminID, RequesterID: otherID})
	assertCode(t, err, domainerrors.CodeBusinessRule)
}

func TestDeactivateLastAdminRejected(t *testing.T) {
	env := newTestEnv(t)
	requesterID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-requester"), "admin", true)
	lastAdminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-lastadmin"), "admin", true)

	if err := env.service.Deactivate(context.Background(), usecase.DeactivateInput{
		TargetID: requesterID, RequesterID: lastAdminID,
	}); err != nil {
		t.Fatalf("setup: deactivate failed: %v", err)
	}

	otherID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-other"), "vendedora", true)
	err := env.service.Deactivate(context.Background(), usecase.DeactivateInput{TargetID: lastAdminID, RequesterID: otherID})
	assertCode(t, err, domainerrors.CodeBusinessRule)
}

func TestActivateSuccess(t *testing.T) {
	env := newTestEnv(t)
	adminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-admin"), "admin", true)
	targetID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-target"), "vendedora", false)

	if err := env.service.Activate(context.Background(), usecase.ActivateInput{TargetID: targetID, RequesterID: adminID}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	user, err := env.service.Get(context.Background(), targetID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !user.IsActive {
		t.Error("expected user to be active")
	}
}

func TestListWithFiltersAndPagination(t *testing.T) {
	env := newTestEnv(t)
	seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-admin"), "admin", true)
	for i := 0; i < 3; i++ {
		seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-vend"+string(rune('a'+i))), "vendedora", true)
	}

	result, err := env.service.List(context.Background(), usecase.ListInput{
		Page: 1, PageSize: 2, Rol: "vendedora", Activo: "true",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 3 {
		t.Errorf("expected total=3 vendedoras, got %d", result.Total)
	}
	if len(result.Items) != 2 {
		t.Errorf("expected 2 items on page 1 (page_size=2), got %d", len(result.Items))
	}
	for _, u := range result.Items {
		if string(u.Rol) != "vendedora" {
			t.Errorf("expected only vendedoras, got rol=%q", u.Rol)
		}
	}
}

func TestUpdatePasswordRequiresCurrentForSelf(t *testing.T) {
	env := newTestEnv(t)
	userID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-self"), "vendedora", true)

	err := env.service.UpdatePassword(context.Background(), usecase.UpdatePasswordInput{
		TargetID: userID, PasswordActual: "wrong-current-password", PasswordNueva: "nueva-password-123",
		RequesterID: userID, RequesterIsAdmin: false,
	})
	assertCode(t, err, domainerrors.CodeValidation)
}

func TestUpdatePasswordSelfWithCorrectCurrentSucceeds(t *testing.T) {
	env := newTestEnv(t)
	userID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-self"), "vendedora", true)

	err := env.service.UpdatePassword(context.Background(), usecase.UpdatePasswordInput{
		TargetID: userID, PasswordActual: testPassword, PasswordNueva: "nueva-password-123",
		RequesterID: userID, RequesterIsAdmin: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdatePasswordAdminCanChangeOthersWithoutCurrent(t *testing.T) {
	env := newTestEnv(t)
	adminID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-admin"), "admin", true)
	targetID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-target"), "vendedora", true)

	err := env.service.UpdatePassword(context.Background(), usecase.UpdatePasswordInput{
		TargetID: targetID, PasswordNueva: "nueva-password-123",
		RequesterID: adminID, RequesterIsAdmin: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdatePasswordNonAdminCannotChangeOthers(t *testing.T) {
	env := newTestEnv(t)
	vendedoraID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-vendedora"), "vendedora", true)
	targetID := seedUsuario(t, env.pool, env.sedeID, uniqueCorreo(t, "-target"), "vendedora", true)

	err := env.service.UpdatePassword(context.Background(), usecase.UpdatePasswordInput{
		TargetID: targetID, PasswordNueva: "nueva-password-123",
		RequesterID: vendedoraID, RequesterIsAdmin: false,
	})
	assertCode(t, err, domainerrors.CodeForbidden)
}

func assertCode(t *testing.T, err error, code domainerrors.Code) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	var domainErr *domainerrors.DomainError
	if !errors.As(err, &domainErr) {
		t.Fatalf("expected *DomainError, got %T: %v", err, err)
	}
	if domainErr.Code != code {
		t.Fatalf("expected code %q, got %q (%v)", code, domainErr.Code, domainErr)
	}
}
