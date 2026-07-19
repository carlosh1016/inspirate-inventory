package productos

import (
	"context"

	"github.com/jackc/pgx/v5"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	commonrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/productos"
	stockactual "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
)

// CreateInput is the request payload plus the requester's context. SedeID is
// inherited from the requester's own claims, not client input. Create is
// admin-only, enforced at the router.
type CreateInput struct {
	SedeID      int64
	Nombre      string
	Categoria   string
	Precio      string
	StockMinimo int32
	RequesterID int64
	IP          string
	UserAgent   string
}

// Create registers a new producto and seeds its stock (vitrina + bodega,
// both zero) atomically.
func (s *Service) Create(ctx context.Context, in CreateInput) (generated.GetProductoByIDRow, error) {
	exists, err := s.Productos.ExistsNombreCategoria(ctx, in.SedeID, in.Nombre, in.Categoria, 0)
	if err != nil {
		return generated.GetProductoByIDRow{}, internalErr(err)
	}
	if exists {
		return generated.GetProductoByIDRow{}, domainerrors.NewConflict(
			"Producto en uso",
			"Ya existe un producto con ese nombre en esa categoría para esta sede.",
		)
	}

	var producto generated.Producto
	err = commonrepo.WithTx(ctx, s.Pool, func(tx pgx.Tx) error {
		txProductos := repo.NewPostgres(tx)
		txStock := stockactual.NewPostgres(tx)

		p, err := txProductos.Insert(ctx, in.SedeID, in.Nombre, in.Categoria, in.Precio, in.StockMinimo)
		if err != nil {
			return err
		}
		producto = p

		return txStock.InitializeStock(ctx, in.SedeID, stockactual.TipoItemProducto, p.ID)
	})
	if err != nil {
		return generated.GetProductoByIDRow{}, internalErr(err)
	}

	result, err := s.Productos.GetByID(ctx, producto.ID)
	if err != nil {
		return generated.GetProductoByIDRow{}, internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "producto_creado", in.IP, in.UserAgent, &producto.ID, nil, snapshot(producto))

	return result, nil
}
