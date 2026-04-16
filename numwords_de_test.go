package speechnorm

import "testing"

// TestGermanToWords is a verbatim port of every de cardinal row from
// Humanizer's LocaleNumberTheoryData.cs. Do not add, remove, or modify
// cases without agreeing with the user first.
//
// Skipped: ordinal rows (erster, zweiter, etc.) — not cardinal.
// Skipped: GrammaticalGender rows — our API has no gender parameter.
func TestGermanToWords(t *testing.T) {
	conv, ok := lookup("de")
	if !ok {
		t.Fatal("german converter not registered")
	}
	cases := []struct {
		n    int64
		want string
	}{
		{0, "null"},
		{1, "eins"},
		{2, "zwei"},
		{3, "drei"},
		{4, "vier"},
		{5, "fünf"},
		{10, "zehn"},
		{11, "elf"},
		{12, "zwölf"},
		{19, "neunzehn"},
		{20, "zwanzig"},
		{21, "einundzwanzig"},
		{22, "zweiundzwanzig"},
		{30, "dreißig"},
		{40, "vierzig"},
		{80, "achtzig"},
		{90, "neunzig"},
		{99, "neunundneunzig"},
		{100, "einhundert"},
		{101, "einhunderteins"},
		{110, "einhundertzehn"},
		{111, "einhundertelf"},
		{115, "einhundertfünfzehn"},
		{121, "einhunderteinundzwanzig"},
		{200, "zweihundert"},
		{999, "neunhundertneunundneunzig"},
		{1000, "eintausend"},
		{1001, "eintausendein"},
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
