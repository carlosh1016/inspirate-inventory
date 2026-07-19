package modelosenvase

import (
	"github.com/go-chi/chi/v5"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
)

// Router mounts /modelos-envase/* under r (typically the /api/v1 group).
// List and get are open to admin and vendedora (needed to browse the
// catalog before creating a variante); create, update and delete are
// admin-only.
func (h *Handler) Router(r chi.Router) {
	r.Route("/modelos-envase", func(r chi.Router) {
		r.Use(middleware.Auth(h.jwtManager))
		r.Use(middleware.RequireRole("admin", "vendedora"))

		r.Get("/", h.List)
		r.Get("/{id}", h.Get)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("admin"))
			r.Post("/", h.Create)
			r.Patch("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)
		})
	})
}
