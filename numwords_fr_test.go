package speechnorm

import "testing"

// TestFrenchToWords is a verbatim port of every fr cardinal row from
// Humanizer's LocaleNumberTheoryData.cs. Do not add, remove, or modify
// cases without agreeing with the user first.
//
// fr-BE and fr-CH rows diverge for 80/90 (nonante/octante) — our registry
// only keys on base "fr", so those variant rows are excluded.
// Skipped: ordinal rows (premier, deuxième, etc.) — not cardinal.
// Skipped: GrammaticalGender rows — our API has no gender parameter.
func TestFrenchToWords(t *testing.T) {
	conv, ok := lookup("fr")
	if !ok {
		t.Fatal("french converter not registered")
	}
	cases := []struct {
		n    int64
		want string
	}{
		{0, "zéro"},
		{1, "un"},
		{2, "deux"},
		{3, "trois"},
		{4, "quatre"},
		{5, "cinq"},
		{10, "dix"},
		{11, "onze"},
		{12, "douze"},
		{19, "dix-neuf"},
		{20, "vingt"},
		{21, "vingt et un"},
		{22, "vingt-deux"},
		{30, "trente"},
		{40, "quarante"},
		{80, "quatre-vingts"},
		{90, "quatre-vingt-dix"},
		{99, "quatre-vingt-dix-neuf"},
		{100, "cent"},
		{101, "cent un"},
		{110, "cent dix"},
		{111, "cent onze"},
		{115, "cent quinze"},
		{121, "cent vingt et un"},
		{200, "deux cents"},
		{999, "neuf cent quatre-vingt-dix-neuf"},
		{1000, "mille"},
		{1001, "mille un"},
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
