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

// CreatePago handles POST /api/v1/cuadres-caja/:id/pagos. Admin and
// vendedora — a pago is routine operative activity.
func (h *Handler) CreatePago(w http.ResponseWriter, r *http.Request) {
	cuadreID, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteError(w, r, badRequestBodyErr())
		return
	}

	var req CreatePagoRequest
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

	pago, err := h.service.AddPago(r.Context(), usecase.AddPagoInput{
		CuadreID:  cuadreID,
		UsuarioID: requester.ID,
		Concepto:  req.Concepto,
		Monto:     monto,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteData(w, http.StatusCreated, toPagoCajaResponse(*pago))
}

// DeletePago handles DELETE /api/v1/cuadres-caja/:id/pagos/:pago_id.
// Admin-only, enforced by the router, and only while the cuadre is still
// abierto.
func (h *Handler) DeletePago(w http.ResponseWriter, r *http.Request) {
	cuadreID, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}
	pagoID, ok := parseIDParam(w, r, "pago_id")
	if !ok {
		return
	}

	if err := h.service.DeletePago(r.Context(), cuadreID, pagoID); err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteNoContent(w)
}
