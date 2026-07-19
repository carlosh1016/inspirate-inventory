package movimientos

import (
	"encoding/json"
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/movimientos"
)

// EntradaMercancia handles POST /api/v1/movimientos/entrada-mercancia.
func (h *Handler) EntradaMercancia(w http.ResponseWriter, r *http.Request) {
	var req EntradaMercanciaRequest
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

	items := make([]usecase.EntradaMercanciaItem, len(req.Items))
	for i, it := range req.Items {
		items[i] = usecase.EntradaMercanciaItem{
			TipoItem:  it.TipoItem,
			ItemID:    it.ItemID,
			Ubicacion: it.Ubicacion,
			Cantidad:  it.Cantidad,
			Motivo:    it.Motivo,
		}
	}

	result, err := h.service.EntradaMercancia(r.Context(), usecase.EntradaMercanciaInput{
		SedeID:      requester.SedeID,
		Items:       items,
		RequesterID: requester.ID,
		IP:          middleware.IPFromContext(r.Context()),
		UserAgent:   middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	responses := make([]MovimientoResponse, len(result))
	for i, m := range result {
		responses[i] = toMovimientoResponseFromDomain(m)
	}
	response.WriteData(w, http.StatusCreated, responses)
}
