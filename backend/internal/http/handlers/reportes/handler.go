// Package reportes contains the HTTP handlers for /api/v1/reportes/* — five
// admin-only endpoints that stream downloadable XLSX files.
package reportes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/jwt"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/reportes"
)

// generarTimeout bounds report generation; exceeding it yields HTTP 504.
const generarTimeout = 60 * time.Second

// Handler groups the reportes HTTP handlers and their dependencies.
type Handler struct {
	service    *usecase.Service
	jwtManager jwt.Manager
	loc        *time.Location
}

// NewHandler builds a Handler. loc is the report timezone (America/Bogota),
// used to interpret the fecha query params.
func NewHandler(service *usecase.Service, jwtManager jwt.Manager, loc *time.Location) *Handler {
	return &Handler{service: service, jwtManager: jwtManager, loc: loc}
}

// respondXLSX writes data as a downloadable, non-cacheable XLSX attachment.
func respondXLSX(w http.ResponseWriter, r *http.Request, filename string, data []byte) {
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	if _, err := w.Write(data); err != nil {
		// Headers are already sent; just log.
		_ = err
	}
	_ = r
}

// mapGenerarErr converts a context deadline to a 504 DomainError as a backstop
// (the usecase already maps ctx errors, so this rarely triggers).
func mapGenerarErr(err error) error {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return domainerrors.NewGatewayTimeout(
			"Reporte demorado",
			"La generación del reporte tomó demasiado tiempo. Intenta con un rango más corto.",
		)
	}
	return err
}

// withTimeout wraps the request context with the generation timeout.
func withTimeout(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), generarTimeout)
}

// requester returns the authenticated user, writing a 401 and returning false
// when absent.
func requester(w http.ResponseWriter, r *http.Request) (middleware.AuthenticatedUser, bool) {
	u, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.WriteError(w, r, domainerrors.NewUnauthorized(
			"No autenticado", "Debes iniciar sesión para continuar."))
		return u, false
	}
	return u, true
}
