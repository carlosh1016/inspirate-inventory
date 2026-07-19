// Package metodospago contains pure domain logic for metodos_pago: no I/O,
// no HTTP, no SQL.
package metodospago

import "time"

// MetodoPago is the domain representation of a metodo_pago.
type MetodoPago struct {
	ID        int64
	Nombre    string
	Codigo    string
	Activo    bool
	DeletedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PuedeEliminarseFisicamente reports whether the metodo_pago has no ventas
// referencing it — the only condition under which it may be hard-deleted
// instead of just soft-deleted.
func PuedeEliminarseFisicamente(ventasAsociadas int64) bool {
	return ventasAsociadas == 0
}
