package ventas

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/ventas"
)

// Update handles PATCH /api/v1/ventas/:id. Admin-only, enforced by the
// router. A venta is immutable except observaciones — the body is decoded
// with DisallowUnknownFields so any other field (e.g. {"total": 0}) is
// rejected as an invalid request instead of being silently ignored.
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

	var req UpdateVentaRequest
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
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

	venta, err := h.service.Update(r.Context(), usecase.UpdateInput{
		TargetID:      id,
		Observaciones: req.Observaciones,
		RequesterID:   requester.ID,
		IP:            middleware.IPFromContext(r.Context()),
		UserAgent:     middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteData(w, http.StatusOK, toVentaDetalladaResponse(venta, nil))
}
