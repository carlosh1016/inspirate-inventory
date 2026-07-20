package ventas

import (
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
)

// ResumenHoy handles GET /api/v1/ventas/hoy/resumen. No query params — uses
// the requester's own sede.
func (h *Handler) ResumenHoy(w http.ResponseWriter, r *http.Request) {
	requester, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.WriteError(w, r, domainerrors.NewUnauthorized("No autenticado", "Debes iniciar sesión para continuar."))
		return
	}

	result, err := h.service.ResumenHoy(r.Context(), requester.SedeID)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	porVendedora := make([]ResumenVendedoraResponse, len(result.PorVendedora))
	for i, row := range result.PorVendedora {
		porVendedora[i] = ResumenVendedoraResponse{
			UsuarioID:      row.UsuarioID,
			NombreCompleto: row.NombreCompleto,
			VentasCount:    row.VentasCount,
			Total:          row.Total.String(),
		}
	}

	topFragancias := make([]ResumenFraganciaResponse, len(result.TopFragancias))
	for i, row := range result.TopFragancias {
		topFragancias[i] = ResumenFraganciaResponse{
			ID:              row.ID,
			NombreComercial: row.NombreComercial,
			GramosVendidos:  row.GramosVendidos.String(),
			MontoVendido:    row.MontoVendido.String(),
		}
	}

	response.WriteData(w, http.StatusOK, ResumenHoyResponse{
		Fecha:          result.Fecha,
		VentasCount:    result.Resumen.VentasCount,
		TotalDia:       result.Resumen.TotalDia.String(),
		DescuentoTotal: result.Resumen.DescuentoTotal.String(),
		PorMetodoPago: ResumenPorMetodoPagoResponse{
			Efectivo:      result.Resumen.TotalEfectivo.String(),
			Nequi:         result.Resumen.TotalNequi.String(),
			Daviplata:     result.Resumen.TotalDaviplata.String(),
			Transferencia: result.Resumen.TotalTransferencia.String(),
			Otros:         result.Resumen.TotalOtros.String(),
		},
		PorVendedora:  porVendedora,
		TopFragancias: topFragancias,
	})
}
