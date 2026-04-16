package speechnorm

import "testing"

// TestPortugueseToWords_ScaleBoundaries covers scale values not in
// Humanizer's LocaleNumberTheoryData.cs. Expected values generated
// by Humanizer itself via Tools/humanizer-reference on 2026-04-14.
// Do NOT modify expected strings without regenerating from Humanizer.
//
// Note: Humanizer's Portuguese converter throws NotImplementedException for
// values >= 1_000_000_000_000, so that case is excluded from this table.
func TestPortugueseToWords_ScaleBoundaries(t *testing.T) {
	conv, _ := lookup("pt")
	cases := []struct {
		n    int64
		want string
	}{
		{1000, "mil"},
		{1001, "mil e um"},
		{1100, "mil e cem"},
		{2000, "dois mil"},
		{9999, "nove mil novecentos e noventa e nove"},
		{10_000, "dez mil"},
		{10_001, "dez mil e um"},
		{99_999, "noventa e nove mil novecentos e noventa e nove"},
		{100_000, "cem mil"},
		{100_001, "cem mil e um"},
		{999_999, "novecentos e noventa e nove mil novecentos e noventa e nove"},
		{1_000_000, "um milhão"},
		{1_000_001, "um milhão e um"},
		{2_000_000, "dois milhões"},
		{9_999_999, "nove milhões novecentos e noventa e nove mil novecentos e noventa e nove"},
		{10_000_000, "dez milhões"},
		{100_000_000, "cem milhões"},
		{1_000_000_000, "mil milhões"},
		{2_000_000_000, "dois mil milhões"},
		{10_000_000_000, "dez mil milhões"},
		{100_000_000_000, "cem mil milhões"},
	}
	for _, c := range cases {
		c := c
		t.Run(intToBase10(c.n), func(t *testing.T) {
			if got := conv.ToWords(c.n); got != c.want {
				t.Errorf("pt ToWords(%d) = %q, want %q", c.n, got, c.want)
			}
		})
	}
}
