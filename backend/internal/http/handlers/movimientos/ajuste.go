package movimientos

import (
	"encoding/json"
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	domainmovimientos "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/movimientos"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/movimientos"
)

// Ajuste handles POST /api/v1/movimientos/ajuste.
func (h *Handler) Ajuste(w http.ResponseWriter, r *http.Request) {
	var req AjusteRequest
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

	m, err := h.service.Ajuste(r.Context(), usecase.AjusteInput{
		SedeID:        requester.SedeID,
		TipoItem:      req.TipoItem,
		ItemID:        req.ItemID,
		Ubicacion:     req.Ubicacion,
		CantidadNueva: req.CantidadNueva,
		Motivo:        req.Motivo,
		RequesterID:   requester.ID,
		IP:            middleware.IPFromContext(r.Context()),
		UserAgent:     middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	writeAjusteResult(w, m)
}

// writeAjusteResult handles the shared response shape for Ajuste and
// Correccion: m is nil when the requested quantity already matched the
// current stock, so no movimiento was created.
func writeAjusteResult(w http.ResponseWriter, m *domainmovimientos.Movimiento) {
	if m == nil {
		response.WriteData(w, http.StatusOK, map[string]string{
			"mensaje": "El stock ya estaba en el valor solicitado; no se generó ningún movimiento.",
		})
		return
	}
	response.WriteData(w, http.StatusCreated, toMovimientoResponseFromDomain(*m))
}
