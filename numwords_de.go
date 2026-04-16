package speechnorm

import "strings"

func init() {
	Register("de", germanConverter{})
}

type germanConverter struct{}

var (
	deUnits = [...]string{
		"null", "eins", "zwei", "drei", "vier", "fünf", "sechs", "sieben",
		"acht", "neun", "zehn", "elf", "zwölf", "dreizehn", "vierzehn",
		"fünfzehn", "sechzehn", "siebzehn", "achtzehn", "neunzehn",
	}
	// Index 2..9 → 20..90
	deTens = [...]string{
		"", "", "zwanzig", "dreißig", "vierzig", "fünfzig",
		"sechzig", "siebzig", "achtzig", "neunzig",
	}
)

func (germanConverter) ToWords(n int64) string {
	if n == 0 {
		return "null"
	}
	if n < 0 {
		return "minus " + germanWords(-n)
	}
	return germanWords(n)
}

func (germanConverter) ToOrdinalWords(n int64) string {
	return germanConverter{}.ToWords(n)
}

// germanWords returns the German cardinal words for n > 0.
func germanWords(n int64) string {
	var parts []string

	if n >= 1_000_000_000_000 {
		billionen := n / 1_000_000_000_000
		n %= 1_000_000_000_000
		prefix := germanMillionPrefix(billionen)
		if billionen == 1 {
			parts = append(parts, prefix+" Billion")
		} else {
			parts = append(parts, prefix+" Billionen")
		}
	}

	if n >= 1_000_000_000 {
		milliarden := n / 1_000_000_000
		n %= 1_000_000_000
		prefix := germanMillionPrefix(milliarden)
		if milliarden == 1 {
			parts = append(parts, prefix+" Milliarde")
		} else {
			parts = append(parts, prefix+" Milliarden")
		}
	}

	if n >= 1_000_000 {
		millionen := n / 1_000_000
		n %= 1_000_000
		prefix := germanMillionPrefix(millionen)
		if millionen == 1 {
			parts = append(parts, prefix+" Million")
		} else {
			parts = append(parts, prefix+" Millionen")
		}
	}

	// Below one million everything concatenates without spaces.
	// compound=true when there are already higher-scale parts: the sub-million
	// chunk is a suffix, so "1" must render as "ein" not standalone "eins".
	if n > 0 {
		parts = append(parts, germanUnderMillion(n, len(parts) > 0))
	}

	return strings.Join(parts, " ")
}

// germanMillionPrefix returns the multiplier word for million/milliarde groups.
// Uses feminine "eine" for 1, otherwise the normal cardinal form.
func germanMillionPrefix(n int64) string {
	if n == 1 {
		return "eine"
	}
	// compound=true: the multiplier word is always embedded in a larger compound.
	return germanUnderMillion(n, true)
}

// germanUnderMillion formats 1..999_999 as a single concatenated string
// (no internal spaces — German compounds everything below one million).
// compound controls whether a trailing "1" renders as "ein" (suffix context)
// or "eins" (standalone context).
func germanUnderMillion(n int64, compound bool) string {
	var b strings.Builder

	if n >= 1_000 {
		thousands := n / 1_000
		n %= 1_000
		b.WriteString(germanUnderThousand(thousands, true))
		b.WriteString("tausend")
		// After "tausend" the remainder is always a suffix — compound=true.
		compound = true
	}

	if n > 0 {
		b.WriteString(germanUnderThousand(n, compound))
	}

	return b.String()
}

// germanUnderThousand formats 1..999.
// compound=true means this chunk is a multiplier (hundred or thousand prefix),
// so "1" renders as "ein" rather than the standalone "eins".
func germanUnderThousand(n int64, compound bool) string {
	if n >= 100 {
		hundredPrefix := germanUnitStem(n/100) + "hundert"
		rem := n % 100
		if rem == 0 {
			return hundredPrefix
		}
		return hundredPrefix + germanUnderHundred(rem, false)
	}
	return germanUnderHundred(n, compound)
}

// germanUnderHundred formats 1..99.
// compound=true suppresses the trailing -s on "eins" → "ein".
func germanUnderHundred(n int64, compound bool) string {
	if n < 20 {
		if n == 1 && compound {
			return "ein"
		}
		return deUnits[n]
	}

	tens := n / 10
	unit := n % 10
	if unit == 0 {
		return deTens[tens]
	}
	// units + "und" + tens — unit 1 uses "ein" in the compound
	stem := germanUnitStem(unit)
	return stem + "und" + deTens[tens]
}

// germanUnitStem returns the compounding stem for units 1..9 used as prefixes
// in tens compounds and as hundred/thousand multipliers.
// For 1 this is "ein"; for all others it is the same as deUnits.
func germanUnitStem(n int64) string {
	if n == 1 {
		return "ein"
	}
	return deUnits[n]
}
