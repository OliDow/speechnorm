package speechnorm

import "testing"

// TestPortugueseToWords is a verbatim port of every pt cardinal row from
// Humanizer's LocaleNumberTheoryData.cs. Do not add, remove, or modify
// cases without agreeing with the user first.
//
// pt-BR rows diverge from base pt at 19 ("dezenove" vs "dezanove"). Our
// registry keys on base "pt" (European Portuguese), so pt-BR variant rows
// are excluded.
// Skipped: ordinal rows — not cardinal.
// Skipped: GrammaticalGender / WordForm rows — our API has no such parameters.
func TestPortugueseToWords(t *testing.T) {
	conv, ok := lookup("pt")
	if !ok {
		t.Fatal("portuguese converter not registered")
	}
	cases := []struct {
		n    int64
		want string
	}{
		{0, "zero"},
		{1, "um"},
		{2, "dois"},
		{3, "três"},
		{4, "quatro"},
		{5, "cinco"},
		{10, "dez"},
		{11, "onze"},
		{12, "doze"},
		{19, "dezanove"},
		{20, "vinte"},
		{21, "vinte e um"},
		{22, "vinte e dois"},
		{30, "trinta"},
		{40, "quarenta"},
		{80, "oitenta"},
		{90, "noventa"},
		{99, "noventa e nove"},
		{100, "cem"},
		{101, "cento e um"},
		{110, "cento e dez"},
		{111, "cento e onze"},
		{115, "cento e quinze"},
		{121, "cento e vinte e um"},
		{200, "duzentos"},
		{999, "novecentos e noventa e nove"},
		{1000, "mil"},
		{1001, "mil e um"},
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
