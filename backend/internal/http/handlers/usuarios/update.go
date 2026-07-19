package usuarios

import (
	"encoding/json"
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/usuarios"
)

// Update handles PATCH /api/v1/usuarios/:id.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r)
	if !ok {
		return
	}

	var req UpdateUsuarioRequest
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

	usuario, err := h.service.Update(r.Context(), usecase.UpdateInput{
		TargetID:       id,
		NombreCompleto: req.NombreCompleto,
		Correo:         req.Correo,
		Rol:            req.Rol,
		RequesterID:    requester.ID,
		IP:             middleware.IPFromContext(r.Context()),
		UserAgent:      middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteData(w, http.StatusOK, toUsuarioResponse(usuario))
}
