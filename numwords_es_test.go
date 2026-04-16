package speechnorm

import "testing"

// TestSpanishToWords is a verbatim port of every es cardinal row from
// Humanizer's LocaleNumberTheoryData.cs. Do not add, remove, or modify
// cases without agreeing with the user first.
//
// es-MX and es-ES rows match the base "es" rows, so no variant exclusions apply.
// Skipped: ordinal rows — not cardinal.
// Skipped: GrammaticalGender rows — our API has no gender parameter.
func TestSpanishToWords(t *testing.T) {
	conv, ok := lookup("es")
	if !ok {
		t.Fatal("spanish converter not registered")
	}
	cases := []struct {
		n    int64
		want string
	}{
		{0, "cero"},
		{1, "uno"},
		{2, "dos"},
		{3, "tres"},
		{4, "cuatro"},
		{5, "cinco"},
		{10, "diez"},
		{11, "once"},
		{12, "doce"},
		{19, "diecinueve"},
		{20, "veinte"},
		{21, "veintiuno"},
		{22, "veintidós"},
		{30, "treinta"},
		{40, "cuarenta"},
		{80, "ochenta"},
		{90, "noventa"},
		{99, "noventa y nueve"},
		{100, "cien"},
		{101, "ciento uno"},
		{110, "ciento diez"},
		{111, "ciento once"},
		{115, "ciento quince"},
		{121, "ciento veintiuno"},
		{200, "doscientos"},
		{999, "novecientos noventa y nueve"},
		{1000, "mil"},
		{1001, "mil uno"},
	}
	for _, c := range cases {
		c := c
		t.Run(intToBase10(c.n), func(t *testing.T) {
			if got := conv.ToWords(c.n); got != c.want {
				t.Errorf("es ToWords(%d) = %q, want %q", c.n, got, c.want)
			}
		})
	}
}
