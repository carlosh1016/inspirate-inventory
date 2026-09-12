package fragancias

import "testing"

func TestParseCodigoFragancia(t *testing.T) {
	cases := []struct {
		in     string
		genero string
		numero int
		ok     bool
	}{
		{"F-014", "femenina", 14, true},
		{"f014", "femenina", 14, true},
		{"F14", "femenina", 14, true},
		{"M-3", "masculina", 3, true},
		{"m3", "masculina", 3, true},
		{"  F-007  ", "femenina", 7, true},
		{"Chanel No. 5", "", 0, false},
		{"", "", 0, false},
		{"F-", "", 0, false},
		{"X-5", "", 0, false},
	}
	for _, c := range cases {
		genero, numero, ok := parseCodigoFragancia(c.in)
		if ok != c.ok || genero != c.genero || numero != c.numero {
			t.Errorf("parseCodigoFragancia(%q) = (%q, %d, %v), want (%q, %d, %v)",
				c.in, genero, numero, ok, c.genero, c.numero, c.ok)
		}
	}
}
