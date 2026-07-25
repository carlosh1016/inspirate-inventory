package reportes

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// Marker types let WriteRow pick the right Excel number format for a value
// whose Go type alone is ambiguous (a decimal could be money, grams, a unit
// count or a percentage; a time could be a full timestamp, a date or a clock
// time). The report code wraps values in these to choose the presentation.
type (
	// Gramos renders a decimal with 2 decimals (matches NUMERIC(10,2)).
	Gramos decimal.Decimal
	// Numero renders a decimal as an integer count with a thousands separator.
	Numero decimal.Decimal
	// Porcentaje renders a decimal as a percentage (e.g. "10%").
	Porcentaje decimal.Decimal
	// Fecha renders a timestamptz instant as a date only (dd/mm/yyyy),
	// converted to the report timezone.
	Fecha time.Time
	// FechaCal renders a calendar date (from a DATE column, already timezone-less)
	// as a date only (dd/mm/yyyy), using its Y/M/D verbatim without timezone
	// conversion — a DATE has no instant, so converting it would shift the day.
	FechaCal time.Time
	// Hora renders a time as a clock time only (hh:mm).
	Hora time.Time
)

// FormatConsecutivoVenta renders a venta id as the human consecutive shown to
// users, e.g. 567 -> "V-000567".
func FormatConsecutivoVenta(id int64) string {
	return fmt.Sprintf("V-%06d", id)
}

// FormatDuracion renders a duration as total "HH:MM:SS", where hours can exceed
// 24 (e.g. 125h30m -> "125:30:00"). Used for aggregated worked hours.
func FormatDuracion(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	total := int64(d / time.Second)
	h := total / 3600
	m := (total % 3600) / 60
	s := total % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

// decimalToFloat converts a decimal to float64 for excelize (which has no
// native decimal type). Safe for COP amounts (integers up to millions) and
// grams limited to 2 decimals.
func decimalToFloat(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}
