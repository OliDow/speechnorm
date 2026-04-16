package speechnorm

import "testing"

// TestFrenchToWords_ScaleBoundaries covers scale values not in
// Humanizer's LocaleNumberTheoryData.cs. Expected values generated
// by Humanizer itself via Tools/humanizer-reference on 2026-04-14.
// Do NOT modify expected strings without regenerating from Humanizer.
func TestFrenchToWords_ScaleBoundaries(t *testing.T) {
	conv, _ := lookup("fr")
	cases := []struct {
		n    int64
		want string
	}{
		{1000, "mille"},
		{1001, "mille un"},
		{1100, "mille cent"},
		{2000, "deux mille"},
		{9999, "neuf mille neuf cent quatre-vingt-dix-neuf"},
		{10_000, "dix mille"},
		{10_001, "dix mille un"},
		{99_999, "quatre-vingt-dix-neuf mille neuf cent quatre-vingt-dix-neuf"},
		{100_000, "cent mille"},
		{100_001, "cent mille un"},
		{999_999, "neuf cent quatre-vingt-dix-neuf mille neuf cent quatre-vingt-dix-neuf"},
		{1_000_000, "un million"},
		{1_000_001, "un million un"},
		{2_000_000, "deux millions"},
		{9_999_999, "neuf millions neuf cent quatre-vingt-dix-neuf mille neuf cent quatre-vingt-dix-neuf"},
		{10_000_000, "dix millions"},
		{100_000_000, "cent millions"},
		{1_000_000_000, "un milliard"},
		{2_000_000_000, "deux milliards"},
		{10_000_000_000, "dix milliards"},
		{100_000_000_000, "cent milliards"},
		{1_000_000_000_000, "un billion"},
	}
	for _, c := range cases {
		c := c
		t.Run(intToBase10(c.n), func(t *testing.T) {
			if got := conv.ToWords(c.n); got != c.want {
				t.Errorf("fr ToWords(%d) = %q, want %q", c.n, got, c.want)
			}
		})
	}
}
