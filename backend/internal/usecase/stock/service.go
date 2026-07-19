// Package stock implements the unified stock view usecase (M8):
// orchestration between repository/stock_actual and nothing else — this
// module is read-only, no auditing, no transactions.
package stock

import (
	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
)

// Service groups the stock usecases behind one set of dependencies.
type Service struct {
	StockActual repo.Repository
}

// NewService builds a Service with all its dependencies.
func NewService(stockActualRepo repo.Repository) *Service {
	return &Service{StockActual: stockActualRepo}
}

func internalErr(err error) error {
	return domainerrors.NewInternal("Error interno", "Ocurrió un error inesperado. Intenta de nuevo más tarde.", err)
}
