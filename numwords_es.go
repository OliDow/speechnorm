package speechnorm

import "strings"

func init() {
	Register("es", spanishConverter{})
}

type spanishConverter struct{}

var (
	// esUnits covers 0..29 — indices 16..19 and 21..29 are contracted forms.
	esUnits = [...]string{
		"cero", "uno", "dos", "tres", "cuatro", "cinco", "seis", "siete",
		"ocho", "nueve", "diez", "once", "doce", "trece", "catorce", "quince",
		"dieciséis", "diecisiete", "dieciocho", "diecinueve", "veinte",
		"veintiuno", "veintidós", "veintitrés", "veinticuatro", "veinticinco",
		"veintiséis", "veintisiete", "veintiocho", "veintinueve",
	}
	// esTens covers indices 3..9 (30..90).
	esTens = [...]string{
		"", "", "", "treinta", "cuarenta", "cincuenta", "sesenta",
		"setenta", "ochenta", "noventa",
	}
	// esHundreds covers indices 1..9 (100..900).
	// Index 1 is "ciento" — used only when there is a remainder; "cien" is
	// the exact-100 special case handled separately.
	esHundreds = [...]string{
		"", "ciento", "doscientos", "trescientos", "cuatrocientos",
		"quinientos", "seiscientos", "setecientos", "ochocientos", "novecientos",
	}
)

func (spanishConverter) ToWords(n int64) string {
	if n == 0 {
		return "cero"
	}
	if n < 0 {
		return "menos " + spanishWords(-n)
	}
	return spanishWords(n)
}

func (spanishConverter) ToOrdinalWords(n int64) string {
	return spanishConverter{}.ToWords(n)
}

// spanishWords returns the Spanish cardinal words for n > 0.
func spanishWords(n int64) string {
	var parts []string

	// Trillions (long scale: 1_000_000_000_000 = "un billón").
	if n >= 1_000_000_000_000 {
		billones := n / 1_000_000_000_000
		n %= 1_000_000_000_000
		var chunk string
		if billones == 1 {
			chunk = "un billón"
		} else {
			chunk = spanishWords(billones) + " billones"
		}
		parts = append(parts, chunk)
	}

	// Billions (long scale: 1_000_000_000 = "mil millones", not "un billón").
	if n >= 1_000_000_000 {
		billions := n / 1_000_000_000
		n %= 1_000_000_000
		// "mil millones" for exactly 1_000_000_000; higher values recurse.
		var chunk string
		if billions == 1 {
			chunk = "mil millones"
		} else {
			chunk = spanishWords(billions) + " mil millones"
		}
		parts = append(parts, chunk)
	}

	if n >= 1_000_000 {
		millions := n / 1_000_000
		n %= 1_000_000
		// Millions are nouns: "un millón" (singular), "dos millones" (plural).
		var chunk string
		if millions == 1 {
			chunk = "un millón"
		} else {
			chunk = spanishWords(millions) + " millones"
		}
		parts = append(parts, chunk)
	}

	if n >= 1_000 {
		thousands := n / 1_000
		n %= 1_000
		// "mil" is invariable — never "un mil".
		if thousands == 1 {
			parts = append(parts, "mil")
		} else {
			parts = append(parts, spanishWords(thousands)+" mil")
		}
	}

	if n >= 100 {
		h := n / 100
		n %= 100
		if h == 1 && n == 0 {
			// Exact 100.
			parts = append(parts, "cien")
		} else if n == 0 {
			parts = append(parts, esHundreds[h])
		} else {
			parts = append(parts, esHundreds[h]+" "+spanishUnderHundred(n))
		}
		n = 0
	}

	if n > 0 {
		parts = append(parts, spanishUnderHundred(n))
	}

	return strings.Join(parts, " ")
}

// spanishUnderHundred formats 1..99.
func spanishUnderHundred(n int64) string {
	if n < 30 {
		return esUnits[n]
	}
	tens := n / 10
	unit := n % 10
	if unit == 0 {
		return esTens[tens]
	}
	return esTens[tens] + " y " + esUnits[unit]
}
