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

// Int4 converts *int32 into a nullable pgtype.Int4.
func Int4(v *int32) pgtype.Int4 {
	if v == nil {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: *v, Valid: true}
}
