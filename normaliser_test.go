package speechnorm

import "testing"

func TestNormaliseNumbers_PlainIntegers(t *testing.T) {
	cases := []struct{ in, locale, want string }{
		{"I have 5 cats", "en", "I have five cats"},
		{"population is 1000000", "en", "population is one million"},
		{"3 cats and 2 dogs", "en", "three cats and two dogs"},
		{"0 items", "en", "zero items"},
		{"there are 42 apples", "en", "there are forty-two apples"},
		{"population is 113965001", "en", "population is one hundred and thirteen million nine hundred and sixty-five thousand and one"},
		{"7500000000 people", "en", "seven billion five hundred million people"},
	}
	runNormaliseCases(t, cases)
}

func TestNormaliseNumbers_CommaSeparated(t *testing.T) {
	cases := []struct{ in, locale, want string }{
		{"cost is 1,000", "en", "cost is one thousand"},
		{"revenue was 1,000,000", "en", "revenue was one million"},
	}
	runNormaliseCases(t, cases)
}

func TestNormaliseNumbers_Ordinals(t *testing.T) {
	cases := []struct{ in, locale, want string }{
		{"1st place", "en", "first place"},
		{"2nd place", "en", "second place"},
		{"3rd place", "en", "third place"},
		{"4th place", "en", "fourth place"},
		{"21st century", "en", "twenty-first century"},
	}
	runNormaliseCases(t, cases)
}

func TestNormaliseNumbers_Currency(t *testing.T) {
	cases := []struct{ in, locale, want string }{
		{"$17.00", "en", "seventeen dollars"},
		{"$17.50", "en", "seventeen dollars and fifty cents"},
		{"$1", "en", "one dollar"},
		{"$1.01", "en", "one dollar and one cent"},
		{"€100", "en", "one hundred euros"},
		{"£1,000", "en", "one thousand pounds"},
		{"£1.50", "en", "one pound and fifty pence"},
		{"¥500", "en", "five hundred yen"},
		{"$1.50", "en", "one dollar and fifty cents"},
		{"$5.5", "en", "five dollars and fifty cents"},
		{"$1.50.", "en", "one dollar and fifty cents."},
		{"$13.75", "en", "thirteen dollars and seventy-five cents"},
		{"I paid $5 for 3 items", "en", "I paid five dollars for three items"},
	}
	runNormaliseCases(t, cases)
}

func TestNormaliseNumbers_DecimalsUnchanged(t *testing.T) {
	cases := []struct{ in, locale, want string }{
		{"pi is 3.14", "en", "pi is 3.14"},
		{"version 2.0 released", "en", "version 2.0 released"},
	}
	runNormaliseCases(t, cases)
}

func TestNormaliseNumbers_Passthrough(t *testing.T) {
	cases := []struct{ in, locale, want string }{
		{"no numbers here", "en", "no numbers here"},
		{"hello world", "en", "hello world"},
	}
	runNormaliseCases(t, cases)
}

func TestNormaliseNumbers_EmptyString(t *testing.T) {
	if got := NormaliseNumbers("", "en"); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestNormaliseNumbers_EmptyLocale(t *testing.T) {
	if got := NormaliseNumbers("I have 5 cats", ""); got != "I have 5 cats" {
		t.Errorf("got %q, want unchanged", got)
	}
}

func TestNormaliseNumbers_UnsupportedLocale(t *testing.T) {
	if got := NormaliseNumbers("I have 5 cats", "zz"); got != "I have 5 cats" {
		t.Errorf("got %q, want unchanged", got)
	}
}

func runNormaliseCases(t *testing.T, cases []struct{ in, locale, want string }) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			if got := NormaliseNumbers(c.in, c.locale); got != c.want {
				t.Errorf("NormaliseNumbers(%q, %q) = %q, want %q",
					c.in, c.locale, got, c.want)
			}
		})
	}
}

func TestNormaliseNumbers_MultiLocaleSmall(t *testing.T) {
	cases := []struct{ in, locale, want string }{
		{"5", "fr", "cinq"},
		{"5", "de", "fünf"},
		{"5", "es", "cinco"},
		{"5", "pt", "cinco"},
		{"5", "it", "cinque"},
		{"5", "ar", "خمسة"},
	}
	runNormaliseCases(t, cases)
}

func TestNormaliseNumbers_MultiLocaleLarge(t *testing.T) {
	cases := []struct{ in, locale, want string }{
		{"1000", "fr", "mille"},
		{"1000", "de", "eintausend"},
		{"1000", "es", "mil"},
		{"1000", "ar", "ألف"},
	}
	runNormaliseCases(t, cases)
}
