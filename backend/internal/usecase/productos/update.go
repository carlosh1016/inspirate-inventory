package productos

import (
	"context"
	"errors"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/productos"
)

// RolVendedora identifies the restricted role for Update: only stock_minimo
// is editable through this endpoint for a vendedora, everything else is
// admin-only.
const RolVendedora = "vendedora"

// UpdateInput is the request payload plus the requester's context. A nil
// field means "leave unchanged". RequesterRole gates which fields a
// vendedora may set — nombre, categoria and precio are admin-only.
type UpdateInput struct {
	TargetID      int64
	Nombre        *string
	Categoria     *string
	Precio        *string
	StockMinimo   *int32
	RequesterID   int64
	RequesterRole string
	IP            string
	UserAgent     string
}

// Update applies a partial update to a producto, enforcing (nombre,
// categoria) uniqueness per sede when either changes. A vendedora may only
// set stock_minimo; any other field in the request is rejected.
func (s *Service) Update(ctx context.Context, in UpdateInput) (generated.GetProductoByIDRow, error) {
	if in.RequesterRole == RolVendedora && (in.Nombre != nil || in.Categoria != nil || in.Precio != nil) {
		return generated.GetProductoByIDRow{}, domainerrors.NewForbidden(
			"Operación no permitida",
			"Solo un administrador puede editar nombre, categoría o precio de un producto.",
		)
	}

	before, err := s.Productos.GetByIDIncludingDeleted(ctx, in.TargetID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return generated.GetProductoByIDRow{}, notFoundErr()
		}
		return generated.GetProductoByIDRow{}, internalErr(err)
	}
	if before.DeletedAt.Valid {
		return generated.GetProductoByIDRow{}, notFoundErr()
	}

	if err := s.checkNombreCategoriaCollision(ctx, before, in); err != nil {
		return generated.GetProductoByIDRow{}, err
	}

	updated, err := s.Productos.Update(ctx, in.TargetID, repo.UpdateFields{
		Nombre:      in.Nombre,
		Categoria:   in.Categoria,
		Precio:      in.Precio,
		StockMinimo: in.StockMinimo,
	})
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return generated.GetProductoByIDRow{}, notFoundErr()
		}
		return generated.GetProductoByIDRow{}, internalErr(err)
	}

	result, err := s.Productos.GetByID(ctx, updated.ID)
	if err != nil {
		return generated.GetProductoByIDRow{}, internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "producto_editado", in.IP, in.UserAgent, &updated.ID, snapshot(before), snapshot(updated))

	return result, nil
}

// checkNombreCategoriaCollision reports a conflict error if in changes
// nombre or categoria to a combination already used by another producto in
// the same sede.
func (s *Service) checkNombreCategoriaCollision(ctx context.Context, before generated.Producto, in UpdateInput) error {
	if in.Nombre == nil && in.Categoria == nil {
		return nil
	}

	nombre := before.Nombre
	if in.Nombre != nil {
		nombre = *in.Nombre
	}
	categoria := string(before.Categoria)
	if in.Categoria != nil {
		categoria = *in.Categoria
	}

	exists, err := s.Productos.ExistsNombreCategoria(ctx, before.SedeID, nombre, categoria, in.TargetID)
	if err != nil {
		return internalErr(err)
	}
	if exists {
		return domainerrors.NewConflict(
			"Producto en uso",
			"Ya existe un producto con ese nombre en esa categoría para esta sede.",
		)
	}
	return nil
}
