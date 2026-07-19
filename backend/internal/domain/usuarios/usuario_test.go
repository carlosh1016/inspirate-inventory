package usuarios_test

import (
	"testing"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/domain/usuarios"
)

func TestCanBeDeactivatedByRejectsSelf(t *testing.T) {
	u := usuarios.Usuario{ID: 1, Rol: "vendedora"}
	err := u.CanBeDeactivatedBy(1, 3)
	assertCode(t, err, domainerrors.CodeBusinessRule)
}

func TestCanBeDeactivatedByRejectsLastAdmin(t *testing.T) {
	u := usuarios.Usuario{ID: 2, Rol: "admin"}
	err := u.CanBeDeactivatedBy(1, 1)
	assertCode(t, err, domainerrors.CodeBusinessRule)
}

func TestCanBeDeactivatedByAllowsWithQuorum(t *testing.T) {
	u := usuarios.Usuario{ID: 2, Rol: "admin"}
	if err := u.CanBeDeactivatedBy(1, 2); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCanBeDeactivatedByAllowsNonAdminRegardlessOfQuorum(t *testing.T) {
	u := usuarios.Usuario{ID: 2, Rol: "vendedora"}
	if err := u.CanBeDeactivatedBy(1, 0); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCanBeDeletedByRejectsSelf(t *testing.T) {
	u := usuarios.Usuario{ID: 1, Rol: "vendedora"}
	err := u.CanBeDeletedBy(1, 3)
	assertCode(t, err, domainerrors.CodeBusinessRule)
}

func TestCanBeDeletedByRejectsLastAdmin(t *testing.T) {
	u := usuarios.Usuario{ID: 2, Rol: "admin"}
	err := u.CanBeDeletedBy(1, 1)
	assertCode(t, err, domainerrors.CodeBusinessRule)
}

func TestCanChangeRoleByRejectsOwnRoleChange(t *testing.T) {
	u := usuarios.Usuario{ID: 1, Rol: "admin"}
	err := u.CanChangeRoleBy(1, "vendedora", 5)
	assertCode(t, err, domainerrors.CodeValidation)
}

func TestCanChangeRoleByAllowsOwnNoOpChange(t *testing.T) {
	u := usuarios.Usuario{ID: 1, Rol: "admin"}
	if err := u.CanChangeRoleBy(1, "admin", 1); err != nil {
		t.Fatalf("expected no error for a no-op role change, got %v", err)
	}
}

func TestCanChangeRoleByRejectsDemotingLastAdmin(t *testing.T) {
	u := usuarios.Usuario{ID: 2, Rol: "admin"}
	err := u.CanChangeRoleBy(1, "vendedora", 1)
	assertCode(t, err, domainerrors.CodeBusinessRule)
}

func TestCanChangeRoleByAllowsDemotingWithQuorum(t *testing.T) {
	u := usuarios.Usuario{ID: 2, Rol: "admin"}
	if err := u.CanChangeRoleBy(1, "vendedora", 2); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCanChangeRoleByAllowsPromotingVendedora(t *testing.T) {
	u := usuarios.Usuario{ID: 2, Rol: "vendedora"}
	if err := u.CanChangeRoleBy(1, "admin", 0); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func assertCode(t *testing.T, err error, code domainerrors.Code) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	domainErr, ok := err.(*domainerrors.DomainError)
	if !ok {
		t.Fatalf("expected *DomainError, got %T", err)
	}
	if domainErr.Code != code {
		t.Fatalf("expected code %q, got %q", code, domainErr.Code)
	}
}
