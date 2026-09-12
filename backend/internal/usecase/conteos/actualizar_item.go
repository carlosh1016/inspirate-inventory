package conteos

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"

	domainconteos "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/conteos"
	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	conteosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/conteos"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// ActualizarGramosFisicoInput is the request payload plus the requester's
// context.
type ActualizarGramosFisicoInput struct {
	ConteoID     int64
	ItemID       int64
	GramosFisico decimal.Decimal
	RequesterID  int64
	IP           string
	UserAgent    string
}

type auditActualizarItemSnapshot struct {
	FraganciaID  int64  `json:"fragancia_id"`
	GramosFisico string `json:"gramos_fisico"`
}

// ActualizarGramosFisico records the physically-counted grams for one
// fragancia within an open conteo. Once the conteo is cerrado, no further
// edits are allowed.
func (s *Service) ActualizarGramosFisico(ctx context.Context, in ActualizarGramosFisicoInput) (*domainconteos.ConteoItem, error) {
	if in.GramosFisico.IsNegative() {
		return nil, domainerrors.NewValidation("Cantidad inválida", "Los gramos físicos no pueden ser negativos.", nil)
	}

	conteo, err := s.Conteos.GetByID(ctx, in.ConteoID)
	if err != nil {
		if errors.Is(err, conteosrepo.ErrNotFound) {
			return nil, notFoundErr()
		}
		return nil, internalErr(err)
	}
	if conteo.Estado != generated.EstadoConteoEnumAbierto {
		return nil, domainerrors.NewConflict("Conteo ya cerrado", "Este conteo de inventario ya fue cerrado y no admite más cambios.")
	}

	row, err := s.Conteos.UpdateItemFisico(ctx, in.ItemID, in.ConteoID, in.GramosFisico)
	if err != nil {
		if errors.Is(err, conteosrepo.ErrNotFound) {
			return nil, itemNotFoundErr()
		}
		return nil, internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "conteo_item_actualizado", in.IP, in.UserAgent, &in.ConteoID, nil, auditActualizarItemSnapshot{
		FraganciaID:  row.FraganciaID,
		GramosFisico: in.GramosFisico.String(),
	})

	item := itemFromUpdateRow(row)
	return &item, nil
}
