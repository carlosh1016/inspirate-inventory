package cuadres

import (
	"context"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"

	domaincuadres "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/cuadres"
	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	cuadresrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/cuadres"
)

// defaultFondoBase is used when AbrirInput.FondoBase is nil.
var defaultFondoBase = decimal.NewFromInt(100000)

// AbrirInput is the request payload plus the requester's context.
type AbrirInput struct {
	SedeID      int64
	FondoBase   *decimal.Decimal
	RequesterID int64
	IP          string
	UserAgent   string
}

// Warning is a soft, non-blocking notice attached to AbrirOutput.
type Warning struct {
	Codigo  string `json:"codigo"`
	Mensaje string `json:"mensaje"`
}

// AbrirOutput is the freshly created cuadre plus any soft warnings (e.g. a
// previous day's cuadre was left open).
type AbrirOutput struct {
	Cuadre   *domaincuadres.Cuadre
	Warnings []Warning
}

// auditAbrirSnapshot is what gets JSON-encoded into datos_despues for
// cuadre_abierto.
type auditAbrirSnapshot struct {
	Fecha     string `json:"fecha"`
	FondoBase string `json:"fondo_base"`
}

// Abrir opens today's cuadre for SedeID. Only one cuadre per (sede, fecha)
// may ever exist (enforced by the UNIQUE constraint too) — a cuadre already
// existing for today, whether abierto or cerrado, is a 409.
func (s *Service) Abrir(ctx context.Context, in AbrirInput) (*AbrirOutput, error) {
	fondoBase := defaultFondoBase
	if in.FondoBase != nil {
		if in.FondoBase.IsNegative() {
			return nil, domainerrors.NewValidation("Fondo base inválido", "El fondo base no puede ser negativo.", nil)
		}
		fondoBase = *in.FondoBase
	}

	hoy := s.hoy()

	if _, err := s.Cuadres.GetBySedeFecha(ctx, in.SedeID, hoy); err == nil {
		return nil, domainerrors.NewConflict("Cuadre ya existe", "Ya existe un cuadre de caja para hoy.")
	} else if !errors.Is(err, cuadresrepo.ErrNotFound) {
		return nil, internalErr(err)
	}

	var warnings []Warning
	if anterior, err := s.Cuadres.GetAbiertoAnterior(ctx, in.SedeID, hoy); err == nil {
		warnings = append(warnings, Warning{
			Codigo:  "cuadre_anterior_abierto",
			Mensaje: fmt.Sprintf("El cuadre del día %s quedó abierto sin cerrar. Ciérralo cuando puedas.", repo.DateString(anterior.Fecha)),
		})
	} else if !errors.Is(err, cuadresrepo.ErrNotFound) {
		return nil, internalErr(err)
	}

	row, err := s.Cuadres.Insert(ctx, cuadresrepo.InsertParams{
		SedeID:    in.SedeID,
		Fecha:     hoy,
		FondoBase: fondoBase,
	})
	if err != nil {
		return nil, internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "cuadre_abierto", in.IP, in.UserAgent, &row.ID, nil, auditAbrirSnapshot{
		Fecha:     repo.DateString(row.Fecha),
		FondoBase: fondoBase.String(),
	})

	cuadre := toDomainCuadre(row)
	return &AbrirOutput{Cuadre: &cuadre, Warnings: warnings}, nil
}
