// Package sesiones holds the pure entity for a vendedora's work session
// (clock in/out): no I/O, no pgx/sqlc types.
package sesiones

import (
	"time"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/domain/cuadres"
)

// Sesion is one clock-in/clock-out span for a usuario. SalidaAt and
// HorasTrabajadas are nil while the session is open.
type Sesion struct {
	ID              int64
	SedeID          int64
	UsuarioID       int64
	EntradaAt       time.Time
	SalidaAt        *time.Time
	HorasTrabajadas *time.Duration
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Usuario         *cuadres.UsuarioBrief
}
