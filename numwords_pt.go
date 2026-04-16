package speechnorm

import "strings"

func init() {
	Register("pt", portugueseConverter{})
}

type portugueseConverter struct{}

var (
	// ptUnits covers 0..19.
	ptUnits = [...]string{
		"zero", "um", "dois", "três", "quatro", "cinco", "seis", "sete",
		"oito", "nove", "dez", "onze", "doze", "treze", "catorze", "quinze",
		"dezasseis", "dezassete", "dezoito", "dezanove",
	}
	// ptTens covers indices 2..9 → 20..90.
	ptTens = [...]string{
		"", "", "vinte", "trinta", "quarenta", "cinquenta",
		"sessenta", "setenta", "oitenta", "noventa",
	}
	// ptHundreds covers indices 1..9 → 100..900 stems.
	// Index 1 is "cento" — used only when there is a remainder; "cem" is
	// the exact-100 special case handled separately.
	ptHundreds = [...]string{
		"", "cento", "duzentos", "trezentos", "quatrocentos",
		"quinhentos", "seiscentos", "setecentos", "oitocentos", "novecentos",
	}
)

func (portugueseConverter) ToWords(n int64) string {
	if n == 0 {
		return "zero"
	}
	if n < 0 {
		return "menos " + portugueseWords(-n)
	}
	return portugueseWords(n)
}

func (portugueseConverter) ToOrdinalWords(n int64) string {
	return portugueseConverter{}.ToWords(n)
}

// portugueseWords returns the Portuguese cardinal words for n > 0.
func portugueseWords(n int64) string {
	var parts []string

	if n >= 1_000_000 {
		millions := n / 1_000_000
		n %= 1_000_000
		// Millions are nouns: "um milhão" (singular), "dois milhões" (plural).
		var milChunk string
		if millions == 1 {
			milChunk = "um milhão"
		} else {
			milChunk = portugueseWords(millions) + " milhões"
		}
		// Use "e" before the remainder when it is "simple" (< 100 or a round
		// multiple of 100). The same rule applies after thousands.
		if n > 0 && ptNeedsE(n) {
			parts = append(parts, milChunk+" e "+portugueseUnder1M(n))
			n = 0
		} else {
			parts = append(parts, milChunk)
		}
	}

	if n >= 1_000 {
		thousands := n / 1_000
		n %= 1_000
		// "mil" is invariable — never "um mil".
		var chunk string
		if thousands == 1 {
			chunk = "mil"
		} else {
			chunk = portugueseWords(thousands) + " mil"
		}
		// Use "e" before the remainder when it is "simple".
		if n > 0 && ptNeedsE(n) {
			parts = append(parts, chunk+" e "+portugueseUnder1000(n))
			n = 0
		} else {
			parts = append(parts, chunk)
		}
	}

	if n > 0 {
		parts = append(parts, portugueseUnder1000(n))
	}

	return strings.Join(parts, " ")
}

// ptNeedsE reports whether the "e" conjunction should precede a remainder
// value when appending it after a higher-scale part. Portuguese inserts "e"
// when the remainder is less than 100 (e.g. "mil e um") or is an exact
// multiple of 100 with no sub-hundred units (e.g. "mil e cem").
func ptNeedsE(n int64) bool {
	return n < 100 || n%100 == 0
}

// portugueseUnder1M formats 1..999_999 without the million-level conjunction
// logic (called inline after conjunction insertion).
func portugueseUnder1M(n int64) string {
	var parts []string
	if n >= 1_000 {
		thousands := n / 1_000
		n %= 1_000
		var chunk string
		if thousands == 1 {
			chunk = "mil"
		} else {
			chunk = portugueseWords(thousands) + " mil"
		}
		if n > 0 && ptNeedsE(n) {
			parts = append(parts, chunk+" e "+portugueseUnder1000(n))
			n = 0
		} else {
			parts = append(parts, chunk)
		}
	}
	if n > 0 {
		parts = append(parts, portugueseUnder1000(n))
	}
	return strings.Join(parts, " ")
}

// portugueseUnder1000 formats 1..999.
func portugueseUnder1000(n int64) string {
	if n >= 100 {
		h := n / 100
		rem := n % 100
		if h == 1 && rem == 0 {
			return "cem"
		}
		if rem == 0 {
			return ptHundreds[h]
		}
		return ptHundreds[h] + " e " + portugueseUnder100(n%100)
	}
	return portugueseUnder100(n)
}

// portugueseUnder100 formats 1..99.
func portugueseUnder100(n int64) string {
	if n < 20 {
		return ptUnits[n]
	}
	tens := n / 10
	unit := n % 10
	if unit == 0 {
		return ptTens[tens]
	}
	return ptTens[tens] + " e " + ptUnits[unit]
}
