// Package stock contains the HTTP handlers for /api/v1/stock/*.
package stock

import (
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/jwt"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/stock"
)

// Handler groups the stock HTTP handlers and their dependencies.
type Handler struct {
	service    *usecase.Service
	jwtManager jwt.Manager
}

// NewHandler builds a Handler.
func NewHandler(service *usecase.Service, jwtManager jwt.Manager) *Handler {
	return &Handler{service: service, jwtManager: jwtManager}
}
