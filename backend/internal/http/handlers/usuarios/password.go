package usuarios

import (
	"encoding/json"
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/usuarios"
)

// UpdatePassword handles PATCH /api/v1/usuarios/:id/password. This route is
// not gated by RequireRole("admin") in the router — admin or the account
// owner may call it — so authorization is checked here explicitly, in
// addition to the usecase checking it again.
func (h *Handler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r)
	if !ok {
		return
	}

	var req UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, r, badRequestBodyErr())
		return
	}
	if err := h.validator.Validate(req); err != nil {
		response.WriteError(w, r, err)
		return
	}

	requester, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.WriteError(w, r, domainerrors.NewUnauthorized("No autenticado", "Debes iniciar sesión para continuar."))
		return
	}

	if requester.ID != id && requester.Rol != "admin" {
		response.WriteError(w, r, domainerrors.NewForbidden(
			"Acceso denegado", "No tienes permiso para realizar esta acción.",
		))
		return
	}

	err := h.service.UpdatePassword(r.Context(), usecase.UpdatePasswordInput{
		TargetID:         id,
		PasswordActual:   req.PasswordActual,
		PasswordNueva:    req.PasswordNueva,
		RequesterID:      requester.ID,
		RequesterIsAdmin: requester.Rol == "admin",
		IP:               middleware.IPFromContext(r.Context()),
		UserAgent:        middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteNoContent(w)
}
