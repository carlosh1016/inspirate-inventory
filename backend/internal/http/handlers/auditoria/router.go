package auditoria

import (
	"github.com/go-chi/chi/v5"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
)

// Router mounts /auditoria/* under r (typically the /api/v1 group). Read-only
// and admin-only — una vendedora nunca consulta la auditoría.
func (h *Handler) Router(r chi.Router) {
	r.Route("/auditoria", func(r chi.Router) {
		r.Use(middleware.Auth(h.jwtManager))
		r.Use(middleware.RequireRole("admin"))

		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
	})
}
