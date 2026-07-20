package ventas

import (
	"fmt"

	"github.com/shopspring/decimal"

	domainventas "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/ventas"
)

type pricingService struct{}

// NewPricingService builds the default PricingService: precio_unitario
// comes from the envase (modelo_envase) or the producto depending on
// TipoLinea; fragancias never carry their own price. Feromona, when
// present, adds precio_feromona*cantidad to the line subtotal without
// creating a separate line.
func NewPricingService() domainventas.PricingService {
	return pricingService{}
}

func (pricingService) Calculate(input domainventas.PricingInput) (domainventas.PricingResult, error) {
	var precioUnitario decimal.Decimal

	switch input.TipoLinea {
	case domainventas.TipoLineaEnvaseConFragancia:
		precioUnitario = input.ModeloEnvase.PrecioConFragancia
	case domainventas.TipoLineaRecarga:
		precioUnitario = input.ModeloEnvase.PrecioRecarga
	case domainventas.TipoLineaEnvaseSolo:
		precioUnitario = input.ModeloEnvase.PrecioSolo
	case domainventas.TipoLineaProductoOtro:
		precioUnitario = input.Producto.Precio
	default:
		return domainventas.PricingResult{}, fmt.Errorf("tipo_linea desconocido: %s", input.TipoLinea)
	}

	cantidad := decimal.NewFromInt32(input.Cantidad)
	subtotal := precioUnitario.Mul(cantidad)

	var precioFeromona decimal.Decimal
	if input.FeromonaProducto != nil {
		precioFeromona = input.FeromonaProducto.Precio
		subtotal = subtotal.Add(precioFeromona.Mul(cantidad))
	}

	return domainventas.PricingResult{
		PrecioUnitario: precioUnitario,
		PrecioFeromona: precioFeromona,
		SubtotalLinea:  subtotal,
	}, nil
}
