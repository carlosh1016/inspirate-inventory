package conteos

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/shopspring/decimal"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/conteos"
)

// ActualizarItem handles PATCH
// /api/v1/conteos-inventario/:id/items/:item_id. Admin-only, enforced by
// the router. Registers the physically-counted grams for one fragancia.
func (h *Handler) ActualizarItem(w http.ResponseWriter, r *http.Request) {
	conteoID, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}
	itemID, ok := parseIDParam(w, r, "item_id")
	if !ok {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteError(w, r, badRequestBodyErr())
		return
	}

	var req ActualizarItemRequest
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

	gramosFisico, err := decimal.NewFromString(req.GramosFisico)
	if err != nil {
		response.WriteError(w, r, domainerrors.NewValidation("Solicitud inválida", "gramos_fisico debe ser un número válido.", nil))
		return
	}

	item, err := h.service.ActualizarGramosFisico(r.Context(), usecase.ActualizarGramosFisicoInput{
		ConteoID:     conteoID,
		ItemID:       itemID,
		GramosFisico: gramosFisico,
		RequesterID:  requester.ID,
		IP:           middleware.IPFromContext(r.Context()),
		UserAgent:    middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteData(w, http.StatusOK, toConteoItemResponse(*item))
}
