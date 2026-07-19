// Package movimientos contains the HTTP handlers for
// /api/v1/movimientos/*.
package movimientos

import (
	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/jwt"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/validator"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/movimientos"
)

// Handler groups the movimientos HTTP handlers and their dependencies.
type Handler struct {
	service    *usecase.Service
	jwtManager jwt.Manager
	validator  *validator.Validator
}

// NewHandler builds a Handler.
func NewHandler(service *usecase.Service, jwtManager jwt.Manager, v *validator.Validator) *Handler {
	return &Handler{service: service, jwtManager: jwtManager, validator: v}
}

func badRequestBodyErr() error {
	return domainerrors.NewValidation("Solicitud inválida", "El cuerpo de la solicitud no es un JSON válido.", nil)
}
