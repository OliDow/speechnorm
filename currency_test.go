package speechnorm

import "testing"

func TestCurrencyNames(t *testing.T) {
	cases := []struct {
		symbol                                  string
		majorSing, majorPlu, minorSing, minorPlu string
	}{
		{"$", "dollar", "dollars", "cent", "cents"},
		{"€", "euro", "euros", "cent", "cents"},
		{"£", "pound", "pounds", "penny", "pence"},
		{"¥", "yen", "yen", "", ""},
		{"₹", "rupee", "rupees", "paisa", "paise"},
	}
	for _, c := range cases {
		t.Run(c.symbol, func(t *testing.T) {
			ms, mp, ns, np := currencyNames(c.symbol)
			if ms != c.majorSing || mp != c.majorPlu || ns != c.minorSing || np != c.minorPlu {
				t.Errorf("currencyNames(%q) = (%q,%q,%q,%q), want (%q,%q,%q,%q)",
					c.symbol, ms, mp, ns, np, c.majorSing, c.majorPlu, c.minorSing, c.minorPlu)
			}
		})
	}
}

func TestCurrencyNamesUnknownSymbol(t *testing.T) {
	ms, mp, ns, np := currencyNames("¤")
	if ms != "" || mp != "" || ns != "" || np != "" {
		t.Errorf("currencyNames(unknown) = (%q,%q,%q,%q), want all empty",
			ms, mp, ns, np)
	}
}

func TestConvertCurrency(t *testing.T) {
	conv := englishConverter{}
	cases := []struct {
		symbol, integerPart, decimalPart, want string
	}{
		{"$", "17", "", "seventeen dollars"},
		{"$", "17", "00", "seventeen dollars"},
		{"$", "17", "50", "seventeen dollars and fifty cents"},
		{"$", "1", "", "one dollar"},
		{"$", "1", "01", "one dollar and one cent"},
		{"€", "100", "", "one hundred euros"},
		{"£", "1,000", "", "one thousand pounds"},
		{"£", "1", "50", "one pound and fifty pence"},
		{"¥", "500", "", "five hundred yen"},
		{"$", "1", "50", "one dollar and fifty cents"},
		{"$", "5", "5", "five dollars and fifty cents"},
		{"$", "13", "75", "thirteen dollars and seventy-five cents"},
		{"¤", "10", "", "ten"}, // unknown symbol: cardinal only
	}
	for _, c := range cases {
		t.Run(c.symbol+c.integerPart+"_"+c.decimalPart, func(t *testing.T) {
			got := convertCurrency(c.symbol, c.integerPart, c.decimalPart, conv)
			if got != c.want {
				t.Errorf("convertCurrency(%q, %q, %q) = %q, want %q",
					c.symbol, c.integerPart, c.decimalPart, got, c.want)
			}
		})
	}
}
