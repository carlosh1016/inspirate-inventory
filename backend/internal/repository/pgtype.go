// Package repository holds small conversion helpers shared by the
// repository implementations in its subpackages (usuarios, refresh_tokens,
// password_resets, ...).
package repository

import (
	"net/netip"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

// Timestamptz converts a time.Time into a valid pgtype.Timestamptz.
func Timestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

// TimePtr converts a pgtype.Timestamptz into *time.Time, nil when NULL.
func TimePtr(ts pgtype.Timestamptz) *time.Time {
	if !ts.Valid {
		return nil
	}
	t := ts.Time
	return &t
}

// TimestamptzPtr converts *time.Time into a nullable pgtype.Timestamptz:
// nil maps to NULL.
func TimestamptzPtr(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

// InetPtr parses an IP string for a nullable INET column. Returns nil (NULL)
// for empty or unparsable input rather than failing the write — the IP is
// diagnostic data, not a business invariant.
func InetPtr(ip string) *netip.Addr {
	if ip == "" {
		return nil
	}
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return nil
	}
	return &addr
}

// Text converts s into a nullable pgtype.Text: empty string maps to NULL.
func Text(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

// StringPtr converts a pgtype.Text into *string, nil when NULL.
func StringPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	s := t.String
	return &s
}

// TextPtr converts *string into a nullable pgtype.Text: nil maps to NULL,
// including a pointer to an empty string (unlike Text, which treats "" as
// NULL too — TextPtr trusts the caller's own nil check).
func TextPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *s, Valid: true}
}

// Int8 converts *int64 into a nullable pgtype.Int8.
func Int8(v *int64) pgtype.Int8 {
	if v == nil {
		return pgtype.Int8{}
	}
	return pgtype.Int8{Int64: *v, Valid: true}
}

// Int8Ptr converts a pgtype.Int8 into *int64, nil when NULL.
func Int8Ptr(v pgtype.Int8) *int64 {
	if !v.Valid {
		return nil
	}
	n := v.Int64
	return &n
}

// Int4 converts *int32 into a nullable pgtype.Int4.
func Int4(v *int32) pgtype.Int4 {
	if v == nil {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: *v, Valid: true}
}

// Date converts a time.Time's calendar date (year/month/day, ignoring
// time-of-day and location) into a valid pgtype.Date.
func Date(t time.Time) pgtype.Date {
	return pgtype.Date{Time: time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), Valid: true}
}

// DatePtr converts *time.Time into a nullable pgtype.Date: nil maps to NULL.
func DatePtr(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{}
	}
	return Date(*t)
}

// DateString formats a pgtype.Date as "YYYY-MM-DD", empty when NULL.
func DateString(d pgtype.Date) string {
	if !d.Valid {
		return ""
	}
	return d.Time.Format("2006-01-02")
}

// IntervalToDuration converts a pgtype.Interval into *time.Duration, nil
// when NULL. Only Days and Microseconds are used: every INTERVAL this
// project reads back was computed in SQL as one timestamptz minus another
// (never constructed with a month component), so Months is always 0.
func IntervalToDuration(iv pgtype.Interval) *time.Duration {
	if !iv.Valid {
		return nil
	}
	d := time.Duration(iv.Days)*24*time.Hour + time.Duration(iv.Microseconds)*time.Microsecond
	return &d
}

// NullDecimal converts *decimal.Decimal into a nullable decimal.NullDecimal.
func NullDecimal(v *decimal.Decimal) decimal.NullDecimal {
	if v == nil {
		return decimal.NullDecimal{}
	}
	return decimal.NullDecimal{Decimal: *v, Valid: true}
}

// NullDecimalPtr converts a decimal.NullDecimal into *decimal.Decimal, nil
// when NULL.
func NullDecimalPtr(v decimal.NullDecimal) *decimal.Decimal {
	if !v.Valid {
		return nil
	}
	d := v.Decimal
	return &d
}
