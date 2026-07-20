package cuadres

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"

	domaincuadres "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/cuadres"
	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	cuadresrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/cuadres"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	pagoscajarepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/pagos_caja"
)

// AddPagoInput is the request payload plus the requester's context.
type AddPagoInput struct {
	CuadreID  int64
	UsuarioID int64
	Concepto  string
	Monto     decimal.Decimal
}

// AddPago registers an operative cash payout against an abierto cuadre.
// Not audited — this is routine operative activity, same criterion as
// individual movimientos in Tanda 3.
func (s *Service) AddPago(ctx context.Context, in AddPagoInput) (*domaincuadres.PagoCaja, error) {
	if !in.Monto.IsPositive() {
		return nil, domainerrors.NewValidation("Monto inválido", "El monto debe ser mayor que 0.", nil)
	}

	cuadre, err := s.Cuadres.GetByID(ctx, in.CuadreID)
	if err != nil {
		if errors.Is(err, cuadresrepo.ErrNotFound) {
			return nil, notFoundErr()
		}
		return nil, internalErr(err)
	}
	if cuadre.Estado != generated.EstadoCuadreEnumAbierto {
		return nil, domainerrors.NewBusinessRule("Cuadre cerrado", "El cuadre de caja está cerrado, no se pueden agregar pagos.")
	}

	row, err := s.Pagos.Insert(ctx, pagoscajarepo.InsertParams{
		CuadreCajaID: in.CuadreID,
		UsuarioID:    in.UsuarioID,
		Concepto:     in.Concepto,
		Monto:        in.Monto,
	})
	if err != nil {
		return nil, internalErr(err)
	}

	nombre := ""
	if usuario, err := s.Usuarios.GetByID(ctx, in.UsuarioID); err == nil {
		nombre = usuario.NombreCompleto
	}

	return &domaincuadres.PagoCaja{
		ID:           row.ID,
		CuadreCajaID: row.CuadreCajaID,
		UsuarioID:    row.UsuarioID,
		Concepto:     row.Concepto,
		Monto:        row.Monto,
		CreatedAt:    row.CreatedAt.Time,
		Usuario:      &domaincuadres.UsuarioBrief{ID: in.UsuarioID, NombreCompleto: nombre},
	}, nil
}
