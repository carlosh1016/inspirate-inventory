// Package repository holds small conversion helpers shared by the
// repository implementations in its subpackages (usuarios, refresh_tokens,
// password_resets, ...).
package repository

import (
	"net/netip"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
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

// Int8 converts *int64 into a nullable pgtype.Int8.
func Int8(v *int64) pgtype.Int8 {
	if v == nil {
		return pgtype.Int8{}
	}
	return pgtype.Int8{Int64: *v, Valid: true}
}
