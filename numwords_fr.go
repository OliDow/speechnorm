package speechnorm

import "strings"

func init() {
	Register("fr", frenchConverter{})
}

type frenchConverter struct{}

var (
	frUnits = [...]string{
		"zéro", "un", "deux", "trois", "quatre", "cinq", "six", "sept",
		"huit", "neuf", "dix", "onze", "douze", "treize", "quatorze",
		"quinze", "seize", "dix-sept", "dix-huit", "dix-neuf",
	}
	frTens = [...]string{
		"", "", "vingt", "trente", "quarante", "cinquante", "soixante",
	}
)

func (frenchConverter) ToWords(n int64) string {
	if n == 0 {
		return "zéro"
	}
	if n < 0 {
		return "moins " + frenchWords(-n)
	}
	return frenchWords(n)
}

func (frenchConverter) ToOrdinalWords(n int64) string {
	// Not exercised by the regex pipeline; return cardinal for interface
	// completeness.
	return frenchConverter{}.ToWords(n)
}

// frenchWords returns the French cardinal words for n > 0.
func frenchWords(n int64) string {
	var parts []string

	if n >= 1_000_000_000_000 {
		billions := n / 1_000_000_000_000
		n %= 1_000_000_000_000
		chunk := frenchWords(billions) + " billion"
		if billions > 1 {
			chunk += "s"
		}
		parts = append(parts, chunk)
	}

	if n >= 1_000_000_000 {
		millions := n / 1_000_000_000
		n %= 1_000_000_000
		chunk := frenchWords(millions) + " milliard"
		if millions > 1 {
			chunk += "s"
		}
		parts = append(parts, chunk)
	}

	if n >= 1_000_000 {
		millions := n / 1_000_000
		n %= 1_000_000
		chunk := frenchWords(millions) + " million"
		if millions > 1 {
			chunk += "s"
		}
		parts = append(parts, chunk)
	}

	if n >= 1_000 {
		thousands := n / 1_000
		n %= 1_000
		// "mille" is invariable — never "un mille", never "milles"
		if thousands == 1 {
			parts = append(parts, "mille")
		} else {
			parts = append(parts, frenchWords(thousands)+" mille")
		}
	}

	if n >= 100 {
		hundreds := n / 100
		n %= 100
		if hundreds == 1 {
			if n == 0 {
				parts = append(parts, "cent")
			} else {
				// "cent" + space + remainder (no -s, not terminal)
				parts = append(parts, "cent "+frenchUnderHundred(n))
			}
		} else {
			prefix := frUnits[hundreds] + " cent"
			if n == 0 {
				// Terminal hundreds take -s: "deux cents", "trois cents", …
				parts = append(parts, prefix+"s")
			} else {
				// Not terminal — no -s
				parts = append(parts, prefix+" "+frenchUnderHundred(n))
			}
		}
		n = 0
	}

	if n > 0 {
		parts = append(parts, frenchUnderHundred(n))
	}

	return strings.Join(parts, " ")
}

// frenchUnderHundred formats 1..99.
func frenchUnderHundred(n int64) string {
	if n < 20 {
		return frUnits[n]
	}

	tens := n / 10
	unit := n % 10

	switch {
	case tens <= 6:
		// 20..69
		if unit == 0 {
			return frTens[tens]
		}
		if unit == 1 {
			return frTens[tens] + " et un"
		}
		return frTens[tens] + "-" + frUnits[unit]

	case tens == 7:
		// 70..79: soixante + 10..19
		inner := unit + 10 // maps to frUnits[10..19]
		if inner == 11 {
			// 71 = soixante et onze
			return "soixante et onze"
		}
		return "soixante-" + frUnits[inner]

	default:
		// tens == 8 or 9: quatre-vingt base
		// 80 = quatre-vingts (terminal -s)
		// 81..89 = quatre-vingt-{unit} (no -s)
		// 90..99 = quatre-vingt-{10..19}
		if tens == 8 {
			if unit == 0 {
				return "quatre-vingts"
			}
			return "quatre-vingt-" + frUnits[unit]
		}
		// tens == 9: 90..99
		inner := unit + 10 // maps to frUnits[10..19]
		return "quatre-vingt-" + frUnits[inner]
	}
}
