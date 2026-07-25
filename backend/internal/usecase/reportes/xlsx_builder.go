package reportes

import (
	"fmt"
	"io"
	"time"
	"unicode/utf8"

	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
)

// rowKind selects the base styling (bold/fill) applied to a row.
type rowKind int

const (
	kindNormal rowKind = iota
	kindHeader
	kindTotal
	kindTitle
)

// numFmt selects the Excel number format applied to a cell.
type numFmt int

const (
	fmtNone numFmt = iota
	fmtMoney
	fmtGramos
	fmtNumero
	fmtPorcentaje
	fmtInteger
	fmtFecha
	fmtHora
	fmtFechaHora
)

func numFmtCode(f numFmt) string {
	switch f {
	case fmtMoney:
		return `"$"#,##0`
	case fmtGramos:
		return `0.00`
	case fmtNumero:
		return `#,##0`
	case fmtPorcentaje:
		return `0.##"%"`
	case fmtInteger:
		return `#,##0`
	case fmtFecha:
		return `dd/mm/yyyy`
	case fmtHora:
		return `hh:mm`
	case fmtFechaHora:
		return `dd/mm/yyyy hh:mm`
	default:
		return ""
	}
}

const (
	colorHeaderFill = "EFEFEF"
	colorTotalFill  = "FFF9DB"
	colorBorder     = "BFBFBF"

	minColWidth = 8.0
	maxColWidth = 50.0
)

// XLSXBuilder wraps an excelize file with a shared style cache and the report
// timezone, so every sheet renders dates in America/Bogota wall time and reuses
// the same styles.
type XLSXBuilder struct {
	file       *excelize.File
	loc        *time.Location
	styleCache map[string]int
}

// NewXLSXBuilder creates an empty workbook. Times are rendered in loc.
func NewXLSXBuilder(loc *time.Location) *XLSXBuilder {
	return &XLSXBuilder{
		file:       excelize.NewFile(),
		loc:        loc,
		styleCache: make(map[string]int),
	}
}

// NewSheet adds a sheet and returns a cursor positioned at row 1.
func (b *XLSXBuilder) NewSheet(name string) *XLSXSheet {
	idx, _ := b.file.NewSheet(name)
	b.file.SetActiveSheet(idx)
	return &XLSXSheet{b: b, name: name, currentRow: 1, colWidths: make(map[int]int)}
}

// DeleteDefaultSheet removes the "Sheet1" excelize creates by default. Call it
// after adding at least one real sheet.
func (b *XLSXBuilder) DeleteDefaultSheet() {
	_ = b.file.DeleteSheet("Sheet1")
}

// Render serializes the workbook to w.
func (b *XLSXBuilder) Render(w io.Writer) error {
	return b.file.Write(w)
}

func (b *XLSXBuilder) styleFor(kind rowKind, f numFmt) int {
	key := fmt.Sprintf("%d:%d", kind, f)
	if id, ok := b.styleCache[key]; ok {
		return id
	}
	st := &excelize.Style{}
	switch kind {
	case kindHeader:
		st.Font = &excelize.Font{Bold: true}
		st.Fill = excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{colorHeaderFill}}
		st.Border = []excelize.Border{{Type: "bottom", Color: colorBorder, Style: 1}}
	case kindTotal:
		st.Font = &excelize.Font{Bold: true}
		st.Fill = excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{colorTotalFill}}
	case kindTitle:
		st.Font = &excelize.Font{Bold: true, Size: 14}
	}
	if code := numFmtCode(f); code != "" {
		st.CustomNumFmt = &code
	}
	id, _ := b.file.NewStyle(st)
	b.styleCache[key] = id
	return id
}

// excelTime normalizes t to the workbook timezone's wall clock, expressed in
// UTC so excelize's serial conversion reflects the local time as displayed.
func (b *XLSXBuilder) excelTime(t time.Time) time.Time {
	tb := t.In(b.loc)
	return time.Date(tb.Year(), tb.Month(), tb.Day(), tb.Hour(), tb.Minute(), tb.Second(), 0, time.UTC)
}

// XLSXSheet is a write cursor over one sheet. It tracks the widest content per
// column so AutoWidth can size them.
type XLSXSheet struct {
	b          *XLSXBuilder
	name       string
	currentRow int
	colWidths  map[int]int
}

// WriteHeaders writes a bold, shaded header row.
func (s *XLSXSheet) WriteHeaders(headers []string) {
	for i, h := range headers {
		col := i + 1
		cell, _ := excelize.CoordinatesToCellName(col, s.currentRow)
		_ = s.b.file.SetCellValue(s.name, cell, h)
		_ = s.b.file.SetCellStyle(s.name, cell, cell, s.b.styleFor(kindHeader, fmtNone))
		s.trackWidth(col, h)
	}
	s.currentRow++
}

// WriteRow writes one data row, dispatching each value to the right format.
func (s *XLSXSheet) WriteRow(values ...interface{}) {
	s.writeValues(kindNormal, values)
}

// WriteTotalRow writes a highlighted total row: label in the first column, then
// the values.
func (s *XLSXSheet) WriteTotalRow(label string, values ...interface{}) {
	all := make([]interface{}, 0, len(values)+1)
	all = append(all, label)
	all = append(all, values...)
	s.writeValues(kindTotal, all)
}

// WriteTitle writes a large bold title merged across mergeAcross columns.
func (s *XLSXSheet) WriteTitle(title string, mergeAcross int) {
	if mergeAcross < 1 {
		mergeAcross = 1
	}
	start, _ := excelize.CoordinatesToCellName(1, s.currentRow)
	end, _ := excelize.CoordinatesToCellName(mergeAcross, s.currentRow)
	_ = s.b.file.SetCellValue(s.name, start, title)
	_ = s.b.file.MergeCell(s.name, start, end)
	_ = s.b.file.SetCellStyle(s.name, start, end, s.b.styleFor(kindTitle, fmtNone))
	s.currentRow++
}

// SkipRow advances the cursor by one, leaving a blank row.
func (s *XLSXSheet) SkipRow() {
	s.currentRow++
}

// AutoWidth sizes each column to its widest tracked content (clamped).
func (s *XLSXSheet) AutoWidth() {
	for col, w := range s.colWidths {
		width := float64(w) + 2
		if width < minColWidth {
			width = minColWidth
		}
		if width > maxColWidth {
			width = maxColWidth
		}
		name, _ := excelize.ColumnNumberToName(col)
		_ = s.b.file.SetColWidth(s.name, name, name, width)
	}
}

func (s *XLSXSheet) writeValues(kind rowKind, values []interface{}) {
	for i, v := range values {
		col := i + 1
		cell, _ := excelize.CoordinatesToCellName(col, s.currentRow)
		cellVal, f, disp := s.resolveCell(v)
		if cellVal != nil {
			_ = s.b.file.SetCellValue(s.name, cell, cellVal)
		}
		if kind != kindNormal || f != fmtNone {
			_ = s.b.file.SetCellStyle(s.name, cell, cell, s.b.styleFor(kind, f))
		}
		s.trackWidth(col, disp)
	}
	s.currentRow++
}

// resolveCell maps a Go value to (excelize cell value, number format, display
// string for width measurement). Numeric/decimal and date/time markers are
// handled by resolveTyped to keep each switch small.
func (s *XLSXSheet) resolveCell(v interface{}) (interface{}, numFmt, string) {
	switch val := v.(type) {
	case nil:
		return nil, fmtNone, ""
	case string:
		return val, fmtNone, val
	case bool:
		if val {
			return "Sí", fmtNone, "Sí"
		}
		return "No", fmtNone, "No"
	case int:
		return val, fmtInteger, thousands(int64(val))
	case int32:
		return val, fmtInteger, thousands(int64(val))
	case int64:
		return val, fmtInteger, thousands(val)
	case *int64:
		if val == nil {
			return nil, fmtNone, ""
		}
		return *val, fmtInteger, thousands(*val)
	case *string:
		if val == nil {
			return nil, fmtNone, ""
		}
		return *val, fmtNone, *val
	case time.Duration:
		disp := FormatDuracion(val)
		return disp, fmtNone, disp
	case *time.Duration:
		if val == nil {
			return nil, fmtNone, ""
		}
		disp := FormatDuracion(*val)
		return disp, fmtNone, disp
	default:
		return s.resolveTyped(v)
	}
}

// resolveTyped handles the decimal and time marker types.
func (s *XLSXSheet) resolveTyped(v interface{}) (interface{}, numFmt, string) {
	switch val := v.(type) {
	case decimal.Decimal:
		return decimalToFloat(val), fmtMoney, "$" + thousands(val.IntPart())
	case Gramos:
		d := decimal.Decimal(val)
		return decimalToFloat(d), fmtGramos, d.StringFixed(2)
	case Numero:
		d := decimal.Decimal(val)
		return decimalToFloat(d), fmtNumero, thousands(d.IntPart())
	case Porcentaje:
		d := decimal.Decimal(val)
		return decimalToFloat(d), fmtPorcentaje, d.String() + "%"
	case Fecha:
		return s.b.excelTime(time.Time(val)), fmtFecha, "dd/mm/yyyy"
	case FechaCal:
		// Calendar date: use its Y/M/D verbatim (no timezone conversion).
		t := time.Time(val)
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), fmtFecha, "dd/mm/yyyy"
	case Hora:
		return s.b.excelTime(time.Time(val)), fmtHora, "hh:mm"
	case time.Time:
		return s.b.excelTime(val), fmtFechaHora, "dd/mm/yyyy hh:mm"
	default:
		disp := fmt.Sprintf("%v", val)
		return disp, fmtNone, disp
	}
}

func (s *XLSXSheet) trackWidth(col int, disp string) {
	if n := utf8.RuneCountInString(disp); n > s.colWidths[col] {
		s.colWidths[col] = n
	}
}

// thousands formats an integer with '.' thousands separators (Colombian style),
// used only to estimate display width.
func thousands(n int64) string {
	neg := n < 0
	if neg {
		n = -n
	}
	digits := fmt.Sprintf("%d", n)
	var out []byte
	for i, c := range []byte(digits) {
		if i > 0 && (len(digits)-i)%3 == 0 {
			out = append(out, '.')
		}
		out = append(out, c)
	}
	if neg {
		return "-" + string(out)
	}
	return string(out)
}
