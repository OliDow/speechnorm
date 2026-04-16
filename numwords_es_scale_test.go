package speechnorm

import "testing"

// TestSpanishToWords_ScaleBoundaries covers scale values not in
// Humanizer's LocaleNumberTheoryData.cs. Expected values generated
// by Humanizer itself via Tools/humanizer-reference on 2026-04-14.
// Do NOT modify expected strings without regenerating from Humanizer.
func TestSpanishToWords_ScaleBoundaries(t *testing.T) {
	conv, _ := lookup("es")
	cases := []struct {
		n    int64
		want string
	}{
		{1000, "mil"},
		{1001, "mil uno"},
		{1100, "mil cien"},
		{2000, "dos mil"},
		{9999, "nueve mil novecientos noventa y nueve"},
		{10_000, "diez mil"},
		{10_001, "diez mil uno"},
		{99_999, "noventa y nueve mil novecientos noventa y nueve"},
		{100_000, "cien mil"},
		{100_001, "cien mil uno"},
		{999_999, "novecientos noventa y nueve mil novecientos noventa y nueve"},
		{1_000_000, "un millón"},
		{1_000_001, "un millón uno"},
		{2_000_000, "dos millones"},
		{9_999_999, "nueve millones novecientos noventa y nueve mil novecientos noventa y nueve"},
		{10_000_000, "diez millones"},
		{100_000_000, "cien millones"},
		{1_000_000_000, "mil millones"},
		{2_000_000_000, "dos mil millones"},
		{10_000_000_000, "diez mil millones"},
		{100_000_000_000, "cien mil millones"},
		{1_000_000_000_000, "un billón"},
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
