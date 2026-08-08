package fragancias

import (
	"github.com/go-chi/chi/v5"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
)

// Router mounts /fragancias/* under r (typically the /api/v1 group). List,
// get, create and update are open to admin and vendedora; delete and
// restaurar require admin.
func (h *Handler) Router(r chi.Router) {
	r.Route("/fragancias", func(r chi.Router) {
		r.Use(middleware.Auth(h.jwtManager))
		r.Use(middleware.RequireRole("admin", "vendedora"))

		r.Get("/", h.List)
		r.Get("/siguiente-numero", h.SiguienteNumero)
		r.Get("/{id}", h.Get)
		r.Post("/", h.Create)
		r.Patch("/{id}", h.Update)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("admin"))
			r.Delete("/{id}", h.Delete)
			r.Post("/{id}/restaurar", h.Restore)
		})
	})
}
