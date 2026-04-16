package speechnorm

import "testing"

// TestItalianToWords is a verbatim port of every it cardinal row from
// Humanizer's LocaleNumberTheoryData.cs. Do not add, remove, or modify
// cases without agreeing with the user first.
//
// Skipped: ordinal rows (primo, secondo, etc.) — not cardinal.
// Skipped: GrammaticalGender / WordForm rows — our API has no such parameters.
// Skipped: addAnd boolean rows — our API has no addAnd parameter.
func TestItalianToWords(t *testing.T) {
	conv, ok := lookup("it")
	if !ok {
		t.Fatal("italian converter not registered")
	}
	cases := []struct {
		n    int64
		want string
	}{
		{0, "zero"},
		{1, "uno"},
		{2, "due"},
		{3, "tre"},
		{4, "quattro"},
		{5, "cinque"},
		{10, "dieci"},
		{11, "undici"},
		{12, "dodici"},
		{19, "diciannove"},
		{20, "venti"},
		{21, "ventuno"},
		{22, "ventidue"},
		{30, "trenta"},
		{40, "quaranta"},
		{80, "ottanta"},
		{90, "novanta"},
		{99, "novantanove"},
		{100, "cento"},
		{101, "centouno"},
		{110, "centodieci"},
		{111, "centoundici"},
		{115, "centoquindici"},
		{121, "centoventuno"},
		{200, "duecento"},
		{999, "novecentonovantanove"},
		{1000, "mille"},
		{1001, "milleuno"},
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
