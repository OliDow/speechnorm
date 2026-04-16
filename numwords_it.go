package speechnorm

func init() {
	Register("it", italianConverter{})
}

type italianConverter struct{}

var (
	itUnits = [...]string{
		"zero", "uno", "due", "tre", "quattro", "cinque", "sei", "sette",
		"otto", "nove", "dieci", "undici", "dodici", "tredici", "quattordici",
		"quindici", "sedici", "diciassette", "diciotto", "diciannove",
	}
	// Index 2..9 → 20..90
	itTens = [...]string{
		"", "", "venti", "trenta", "quaranta", "cinquanta",
		"sessanta", "settanta", "ottanta", "novanta",
	}
)

func (italianConverter) ToWords(n int64) string {
	if n == 0 {
		return "zero"
	}
	if n < 0 {
		return "meno " + italianWords(-n)
	}
	return italianWords(n)
}

func (italianConverter) ToOrdinalWords(n int64) string {
	return italianConverter{}.ToWords(n)
}

// italianWords returns Italian cardinal words for n > 0.
func italianWords(n int64) string {
	if n >= 1_000_000_000 {
		miliardi := n / 1_000_000_000
		rem := n % 1_000_000_000
		var chunk string
		if miliardi == 1 {
			chunk = "un miliardo"
		} else {
			chunk = italianWords(miliardi) + " miliardi"
		}
		if rem == 0 {
			return chunk
		}
		return chunk + " " + italianWords(rem)
	}

	if n >= 1_000_000 {
		millions := n / 1_000_000
		rem := n % 1_000_000
		var chunk string
		if millions == 1 {
			chunk = "un milione"
		} else {
			chunk = italianWords(millions) + " milioni"
		}
		if rem == 0 {
			return chunk
		}
		return chunk + " " + italianWords(rem)
	}

	// Below one million everything is concatenated without spaces.
	return italianSubMillion(n)
}

// italianSubMillion formats 1..999_999 without spaces (thousands and below).
func italianSubMillion(n int64) string {
	if n >= 1_000 {
		thousands := n / 1_000
		rem := n % 1_000
		thousandsWord := itThousandsBlock(thousands)
		if rem == 0 {
			return thousandsWord
		}
		return thousandsWord + italianUnder1000(rem)
	}
	return italianUnder1000(n)
}

// itThousandsBlock returns the thousands-scale word for a given multiplier.
// 1 → "mille", 2 → "duemila", 3 → "tremila", etc.
func itThousandsBlock(thousands int64) string {
	if thousands == 1 {
		return "mille"
	}
	return italianSubMillion(thousands) + "mila"
}

// italianUnder1000 formats 1..999.
func italianUnder1000(n int64) string {
	if n >= 100 {
		return itHundreds(n)
	}
	if n >= 20 {
		return itTensUnits(n)
	}
	return itUnits[n]
}

// itTensUnits formats 20..99 with vowel elision and terminal tré accent.
func itTensUnits(n int64) string {
	tens := n / 10
	unit := n % 10
	if unit == 0 {
		return itTens[tens]
	}
	tensWord := itTens[tens]
	unitWord := itUnits[unit]
	// Vowel elision: drop final vowel of tens word when unit is 1 or 8.
	// All tens words (venti, trenta, quaranta, etc.) end in an ASCII vowel.
	if unit == 1 || unit == 8 {
		tensWord = tensWord[:len(tensWord)-1]
	}
	// Terminal tré accent: 3 becomes tré when it is the final unit of the number.
	if unit == 3 {
		unitWord = "tré"
	}
	return tensWord + unitWord
}

// itHundreds formats 100..999.
func itHundreds(n int64) string {
	h := n / 100
	rem := n % 100
	var prefix string
	if h == 1 {
		prefix = "cento"
	} else {
		// Hundreds prefix uses plain itUnits (tre, not tré — non-terminal).
		prefix = itUnits[h] + "cento"
	}
	if rem == 0 {
		return prefix
	}
	// No vowel elision between hundreds and remainder — concatenate directly.
	if rem < 20 {
		unitWord := itUnits[rem]
		// Terminal tré accent applies here too.
		if rem == 3 {
			unitWord = "tré"
		}
		return prefix + unitWord
	}
	return prefix + itTensUnits(rem)
}
