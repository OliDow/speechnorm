package speechnorm

import "testing"

// TestItalianToWords_ScaleBoundaries covers scale values not in
// Humanizer's LocaleNumberTheoryData.cs. Expected values generated
// by Humanizer itself via Tools/humanizer-reference on 2026-04-14.
// Do NOT modify expected strings without regenerating from Humanizer.
//
// Note: Humanizer's Italian converter throws NotImplementedException for
// values >= 10_000_000_000, so cases above 2_000_000_000 are excluded.
func TestItalianToWords_ScaleBoundaries(t *testing.T) {
	conv, _ := lookup("it")
	cases := []struct {
		n    int64
		want string
	}{
		{1000, "mille"},
		{1001, "milleuno"},
		{1100, "millecento"},
		{2000, "duemila"},
		{9999, "novemilanovecentonovantanove"},
		{10_000, "diecimila"},
		{10_001, "diecimilauno"},
		{99_999, "novantanovemilanovecentonovantanove"},
		{100_000, "centomila"},
		{100_001, "centomilauno"},
		{999_999, "novecentonovantanovemilanovecentonovantanove"},
		{1_000_000, "un milione"},
		{1_000_001, "un milione uno"},
		{2_000_000, "due milioni"},
		{9_999_999, "nove milioni novecentonovantanovemilanovecentonovantanove"},
		{10_000_000, "dieci milioni"},
		{100_000_000, "cento milioni"},
		{1_000_000_000, "un miliardo"},
		{2_000_000_000, "due miliardi"},
	}
	for _, c := range cases {
		c := c
		t.Run(intToBase10(c.n), func(t *testing.T) {
			if got := conv.ToWords(c.n); got != c.want {
				t.Errorf("it ToWords(%d) = %q, want %q", c.n, got, c.want)
			}
		})
	}
}
