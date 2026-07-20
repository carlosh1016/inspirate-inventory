package cuadres

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"

	domaincuadres "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/cuadres"
	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/consignaciones"
	cuadresrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/cuadres"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// AddConsignacionInput is the request payload plus the requester's context.
type AddConsignacionInput struct {
	CuadreID   int64
	UsuarioID  int64
	Monto      decimal.Decimal
	Banco      *string
	Referencia *string
}

// AddConsignacion registers a bank deposit against an abierto cuadre. Not
// audited — routine operative activity, same criterion as AddPago.
func (s *Service) AddConsignacion(ctx context.Context, in AddConsignacionInput) (*domaincuadres.Consignacion, error) {
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
		return nil, domainerrors.NewBusinessRule("Cuadre cerrado", "El cuadre de caja está cerrado, no se pueden agregar consignaciones.")
	}

	row, err := s.Consignaciones.Insert(ctx, consignaciones.InsertParams{
		CuadreCajaID: in.CuadreID,
		UsuarioID:    in.UsuarioID,
		Monto:        in.Monto,
		Banco:        in.Banco,
		Referencia:   in.Referencia,
	})
	if err != nil {
		return nil, internalErr(err)
	}

	nombre := ""
	if usuario, err := s.Usuarios.GetByID(ctx, in.UsuarioID); err == nil {
		nombre = usuario.NombreCompleto
	}

	return &domaincuadres.Consignacion{
		ID:           row.ID,
		CuadreCajaID: row.CuadreCajaID,
		UsuarioID:    row.UsuarioID,
		Monto:        row.Monto,
		Banco:        in.Banco,
		Referencia:   in.Referencia,
		CreatedAt:    row.CreatedAt.Time,
		Usuario:      &domaincuadres.UsuarioBrief{ID: in.UsuarioID, NombreCompleto: nombre},
	}, nil
}
