package cuadres

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/shopspring/decimal"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/cuadres"
)

// CreateConsignacion handles POST /api/v1/cuadres-caja/:id/consignaciones.
// Admin and vendedora — a consignación is routine operative activity.
func (h *Handler) CreateConsignacion(w http.ResponseWriter, r *http.Request) {
	cuadreID, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteError(w, r, badRequestBodyErr())
		return
	}

	var req CreateConsignacionRequest
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

	monto, err := decimal.NewFromString(req.Monto)
	if err != nil {
		response.WriteError(w, r, domainerrors.NewValidation("Solicitud inválida", "monto debe ser un número válido.", nil))
		return
	}

	consignacion, err := h.service.AddConsignacion(r.Context(), usecase.AddConsignacionInput{
		CuadreID:   cuadreID,
		UsuarioID:  requester.ID,
		Monto:      monto,
		Banco:      req.Banco,
		Referencia: req.Referencia,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteData(w, http.StatusCreated, toConsignacionResponse(*consignacion))
}

// DeleteConsignacion handles DELETE
// /api/v1/cuadres-caja/:id/consignaciones/:consignacion_id. Admin-only,
// enforced by the router, and only while the cuadre is still abierto.
func (h *Handler) DeleteConsignacion(w http.ResponseWriter, r *http.Request) {
	cuadreID, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}
	consignacionID, ok := parseIDParam(w, r, "consignacion_id")
	if !ok {
		return
	}

	if err := h.service.DeleteConsignacion(r.Context(), cuadreID, consignacionID); err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteNoContent(w)
}
