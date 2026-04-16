package speechnorm

import "testing"

// TestEnglishToWords is a verbatim port of every en/en-US/en-GB cardinal row
// from Humanizer's NumberToWordsTests.cs (ToWordsInt + ToWordsLong methods)
// and LocaleNumberTheoryData.cs as of github.com/Humanizr/Humanizer main branch.
//
// ToWordsWithoutAnd (useAnd=false) is skipped — our API has no useAnd parameter.
// ToWords_CanSpecifyCultureExplicitly cases for non-en locales are skipped —
// they are covered by their own locale test files.
func TestEnglishToWords(t *testing.T) {
	conv := englishConverter{}
	cases := []struct {
		n    int64
		want string
	}{
		// ---- ToWordsInt cases (Humanizer NumberToWordsTests.cs) ----
		{-1, "minus one"},
		{1, "one"},
		{10, "ten"},
		{11, "eleven"},
		{20, "twenty"},
		{100, "one hundred"},
		{111, "one hundred and eleven"},
		{122, "one hundred and twenty-two"},
		{123, "one hundred and twenty-three"},
		{1000, "one thousand"},
		{1001, "one thousand and one"},
		{1010, "one thousand and ten"},
		{1111, "one thousand one hundred and eleven"},
		{1234, "one thousand two hundred and thirty-four"},
		{3501, "three thousand five hundred and one"},
		{12345, "twelve thousand three hundred and forty-five"},
		{100000, "one hundred thousand"},
		{111111, "one hundred and eleven thousand one hundred and eleven"},
		{123456, "one hundred and twenty-three thousand four hundred and fifty-six"},
		{1000000, "one million"},
		{1111111, "one million one hundred and eleven thousand one hundred and eleven"},
		{1234567, "one million two hundred and thirty-four thousand five hundred and sixty-seven"},
		{10000000, "ten million"},
		{11111111, "eleven million one hundred and eleven thousand one hundred and eleven"},
		{12345678, "twelve million three hundred and forty-five thousand six hundred and seventy-eight"},
		{100000000, "one hundred million"},
		{111111111, "one hundred and eleven million one hundred and eleven thousand one hundred and eleven"},
		{123456789, "one hundred and twenty-three million four hundred and fifty-six thousand seven hundred and eighty-nine"},
		{1000000000, "one billion"},
		{1111111111, "one billion one hundred and eleven million one hundred and eleven thousand one hundred and eleven"},
		{1234567890, "one billion two hundred and thirty-four million five hundred and sixty-seven thousand eight hundred and ninety"},

		// ---- ToWordsLong cases — extends into trillion/quadrillion/quintillion ----
		{11111111111, "eleven billion one hundred and eleven million one hundred and eleven thousand one hundred and eleven"},
		{111111111111, "one hundred and eleven billion one hundred and eleven million one hundred and eleven thousand one hundred and eleven"},
		{1111111111111, "one trillion one hundred and eleven billion one hundred and eleven million one hundred and eleven thousand one hundred and eleven"},
		{11111111111111, "eleven trillion one hundred and eleven billion one hundred and eleven million one hundred and eleven thousand one hundred and eleven"},
		{111111111111111, "one hundred and eleven trillion one hundred and eleven billion one hundred and eleven million one hundred and eleven thousand one hundred and eleven"},
		{1111111111111111, "one quadrillion one hundred and eleven trillion one hundred and eleven billion one hundred and eleven million one hundred and eleven thousand one hundred and eleven"},
		{11111111111111111, "eleven quadrillion one hundred and eleven trillion one hundred and eleven billion one hundred and eleven million one hundred and eleven thousand one hundred and eleven"},
		{111111111111111111, "one hundred and eleven quadrillion one hundred and eleven trillion one hundred and eleven billion one hundred and eleven million one hundred and eleven thousand one hundred and eleven"},
		{1111111111111111111, "one quintillion one hundred and eleven quadrillion one hundred and eleven trillion one hundred and eleven billion one hundred and eleven million one hundred and eleven thousand one hundred and eleven"},

		// ---- LocaleNumberTheoryData.cs: en / en-US / en-GB rows ----
		// (0–1001 range; values already present above are kept for completeness)
		{0, "zero"},
		{2, "two"},
		{3, "three"},
		{4, "four"},
		{5, "five"},
		{12, "twelve"},
		{19, "nineteen"},
		{21, "twenty-one"},
		{22, "twenty-two"},
		{30, "thirty"},
		{40, "forty"},
		{80, "eighty"},
		{90, "ninety"},
		{99, "ninety-nine"},
		{101, "one hundred and one"},
		{110, "one hundred and ten"},
		{115, "one hundred and fifteen"},
		{121, "one hundred and twenty-one"},
		{200, "two hundred"},
		{999, "nine hundred and ninety-nine"},

		// ---- ToWords_CanSpecifyCultureExplicitly en-US row ----
		// (same expected value as generic English)
		// {11, "eleven"} -- covered above
	}
	for _, c := range cases {
		c := c
		t.Run(englishName(c.n), func(t *testing.T) {
			if got := conv.ToWords(c.n); got != c.want {
				t.Errorf("ToWords(%d) = %q, want %q", c.n, got, c.want)
			}
		})
	}
}

// TestEnglishToOrdinalWords is a verbatim port of every English ordinal row
// from Humanizer's NumberToWordsTests.cs (ToOrdinalWords method) and the
// ToOrdinalWords_CanSpecifyCultureExplicitly en-US row.
func TestEnglishToOrdinalWords(t *testing.T) {
	conv := englishConverter{}
	cases := []struct {
		n    int64
		want string
	}{
		{0, "zeroth"},
		{1, "first"},
		{2, "second"},
		{3, "third"},
		{4, "fourth"},
		{5, "fifth"},
		{6, "sixth"},
		{7, "seventh"},
		{8, "eighth"},
		{9, "ninth"},
		{10, "tenth"},
		{11, "eleventh"},
		{12, "twelfth"},
		{13, "thirteenth"},
		{14, "fourteenth"},
		{15, "fifteenth"},
		{16, "sixteenth"},
		{17, "seventeenth"},
		{18, "eighteenth"},
		{19, "nineteenth"},
		{20, "twentieth"},
		{21, "twenty-first"},
		{22, "twenty-second"},
		{30, "thirtieth"},
		{40, "fortieth"},
		{50, "fiftieth"},
		{60, "sixtieth"},
		{70, "seventieth"},
		{80, "eightieth"},
		{90, "ninetieth"},
		{95, "ninety-fifth"},
		{96, "ninety-sixth"},
		// BUG-1: round-scale ordinals drop the leading "one"
		{100, "hundredth"},
		{112, "hundred and twelfth"},
		{120, "hundred and twentieth"},
		{121, "hundred and twenty-first"},
		{1000, "thousandth"},
		{1001, "thousand and first"},
		{1021, "thousand and twenty-first"},
		{10000, "ten thousandth"},
		{10121, "ten thousand one hundred and twenty-first"},
		{100000, "hundred thousandth"},
		{1000000, "millionth"},
		// ToOrdinalWords_CanSpecifyCultureExplicitly en-US row
		// {1021, "en-US", "thousand and twenty-first"} -- already covered above
	}
	for _, c := range cases {
		c := c
		t.Run(englishName(c.n), func(t *testing.T) {
			if got := conv.ToOrdinalWords(c.n); got != c.want {
				t.Errorf("ToOrdinalWords(%d) = %q, want %q", c.n, got, c.want)
			}
		})
	}
}
