// Package usuarios contains pure domain rules for user management: no I/O,
// no HTTP, no SQL. Usecases fetch whatever counts these rules need (e.g.
// active admin count) and pass them in.
package usuarios

import (
	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
)

// Usuario is the domain representation of a system user, carrying just the
// fields these business-rule checks need.
type Usuario struct {
	ID       int64
	Rol      string
	IsActive bool
}

// CanBeDeactivatedBy reports whether requesterID may deactivate this
// usuario: never themselves, and — if this usuario is the last active
// admin — never at all. activeAdminCount is the count before the change.
func (u Usuario) CanBeDeactivatedBy(requesterID int64, activeAdminCount int64) error {
	if u.ID == requesterID {
		return domainerrors.NewBusinessRule(
			"Operación no permitida",
			"No puedes desactivar tu propia cuenta.",
		)
	}
	if u.Rol == "admin" && activeAdminCount <= 1 {
		return domainerrors.NewBusinessRule(
			"Operación no permitida",
			"Debe quedar al menos un administrador activo en el sistema.",
		)
	}
	return nil
}

// CanBeDeletedBy applies the same invariants as CanBeDeactivatedBy —
// deleting is a stronger deactivation.
func (u Usuario) CanBeDeletedBy(requesterID int64, activeAdminCount int64) error {
	if u.ID == requesterID {
		return domainerrors.NewBusinessRule(
			"Operación no permitida",
			"No puedes eliminar tu propia cuenta.",
		)
	}
	if u.Rol == "admin" && activeAdminCount <= 1 {
		return domainerrors.NewBusinessRule(
			"Operación no permitida",
			"Debe quedar al menos un administrador activo en el sistema.",
		)
	}
	return nil
}

// CanChangeRoleBy reports whether requesterID may change this usuario's rol
// to newRol: an admin can never change their own rol (self-lockout guard,
// checked first regardless of quorum), and demoting the last active admin
// is never allowed. activeAdminCount is the count before the change.
func (u Usuario) CanChangeRoleBy(requesterID int64, newRol string, activeAdminCount int64) error {
	if u.ID == requesterID && newRol != u.Rol {
		return domainerrors.NewValidation(
			"Operación no permitida",
			"No puedes cambiar tu propio rol.",
			nil,
		)
	}
	if u.Rol == "admin" && newRol != "admin" && activeAdminCount <= 1 {
		return domainerrors.NewBusinessRule(
			"Operación no permitida",
			"Debe quedar al menos un administrador activo en el sistema.",
		)
	}
	return nil
}
