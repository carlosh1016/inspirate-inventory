package sesiones

import (
	"encoding/json"
	"io"
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/sesiones"
)

// Update handles PATCH /api/v1/sesiones-laborales/:id. Admin-only,
// enforced by the router — manual correction of entrada_at/salida_at;
// horas_trabajadas is recalculated automatically.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r)
	if !ok {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteError(w, r, badRequestBodyErr())
		return
	}

	var req UpdateSesionRequest
	if err := json.Unmarshal(body, &req); err != nil {
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

	sesion, err := h.service.Update(r.Context(), usecase.UpdateInput{
		TargetID:    id,
		EntradaAt:   req.EntradaAt,
		SalidaAt:    req.SalidaAt,
		RequesterID: requester.ID,
		IP:          middleware.IPFromContext(r.Context()),
		UserAgent:   middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteData(w, http.StatusOK, toSesionResponse(*sesion))
}
