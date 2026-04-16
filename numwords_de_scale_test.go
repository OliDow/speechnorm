package speechnorm

import "testing"

// TestGermanToWords_ScaleBoundaries covers scale values not in
// Humanizer's LocaleNumberTheoryData.cs. Expected values generated
// by Humanizer itself via Tools/humanizer-reference on 2026-04-14.
// Do NOT modify expected strings without regenerating from Humanizer.
func TestGermanToWords_ScaleBoundaries(t *testing.T) {
	conv, _ := lookup("de")
	cases := []struct {
		n    int64
		want string
	}{
		{1000, "eintausend"},
		{1001, "eintausendein"},
		{1100, "eintausendeinhundert"},
		{2000, "zweitausend"},
		{9999, "neuntausendneunhundertneunundneunzig"},
		{10_000, "zehntausend"},
		{10_001, "zehntausendein"},
		{99_999, "neunundneunzigtausendneunhundertneunundneunzig"},
		{100_000, "einhunderttausend"},
		{100_001, "einhunderttausendein"},
		{999_999, "neunhundertneunundneunzigtausendneunhundertneunundneunzig"},
		{1_000_000, "eine Million"},
		{1_000_001, "eine Million ein"},
		{2_000_000, "zwei Millionen"},
		{9_999_999, "neun Millionen neunhundertneunundneunzigtausendneunhundertneunundneunzig"},
		{10_000_000, "zehn Millionen"},
		{100_000_000, "einhundert Millionen"},
		{1_000_000_000, "eine Milliarde"},
		{2_000_000_000, "zwei Milliarden"},
		{10_000_000_000, "zehn Milliarden"},
		{100_000_000_000, "einhundert Milliarden"},
		{1_000_000_000_000, "eine Billion"},
	}
	for _, c := range cases {
		c := c
		t.Run(intToBase10(c.n), func(t *testing.T) {
			if got := conv.ToWords(c.n); got != c.want {
				t.Errorf("de ToWords(%d) = %q, want %q", c.n, got, c.want)
			}
		})
	}
}
