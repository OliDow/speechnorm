package speechnorm

import "strings"

func init() {
	Register("en", englishConverter{})
}

type englishConverter struct{}

var (
	enUnits = [...]string{
		"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
		"ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen",
		"seventeen", "eighteen", "nineteen",
	}
	enTens = [...]string{
		"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety",
	}
)

// enScales defines the short-scale groups in descending order, up to int64 max
// (~9.2 quintillion). Each entry covers values from value up to the next entry.
var enScales = []struct {
	value int64
	name  string
}{
	{1_000_000_000_000_000_000, "quintillion"},
	{1_000_000_000_000_000, "quadrillion"},
	{1_000_000_000_000, "trillion"},
	{1_000_000_000, "billion"},
	{1_000_000, "million"},
	{1_000, "thousand"},
}

func (englishConverter) ToWords(n int64) string {
	if n < 0 {
		return "minus " + englishWords(-n)
	}
	return englishWords(n)
}

func englishWords(n int64) string {
	if n == 0 {
		return "zero"
	}

	var parts []string
	for _, s := range enScales {
		if n >= s.value {
			parts = append(parts, englishUnderThousand(n/s.value)+" "+s.name)
			n %= s.value
		}
	}
	if n > 0 {
		// "and" only when we already have a higher-order group and the
		// remainder is under 100 OR when the remainder has no hundreds
		// group of its own.
		if len(parts) > 0 && n < 100 {
			parts = append(parts, "and "+englishUnderThousand(n))
		} else {
			parts = append(parts, englishUnderThousand(n))
		}
	}
	return strings.Join(parts, " ")
}

// englishUnderThousand formats 1..999 with internal "and".
func englishUnderThousand(n int64) string {
	if n < 20 {
		return enUnits[n]
	}
	if n < 100 {
		if n%10 == 0 {
			return enTens[n/10]
		}
		return enTens[n/10] + "-" + enUnits[n%10]
	}
	hundreds := enUnits[n/100] + " hundred"
	rem := n % 100
	if rem == 0 {
		return hundreds
	}
	return hundreds + " and " + englishUnderThousand(rem)
}

var enOrdinalsIrregular = map[int64]string{
	0: "zeroth", 1: "first", 2: "second", 3: "third", 4: "fourth",
	5: "fifth", 6: "sixth", 7: "seventh", 8: "eighth", 9: "ninth",
	10: "tenth", 11: "eleventh", 12: "twelfth",
}

// enTensOrdinal: twentieth, thirtieth … — used when n%10 == 0 and n in [20,90].
var enTensOrdinal = map[int64]string{
	20: "twentieth", 30: "thirtieth", 40: "fortieth", 50: "fiftieth",
	60: "sixtieth", 70: "seventieth", 80: "eightieth", 90: "ninetieth",
}

// enLastWordOrdinal replaces the last word of a phrase with its ordinal form.
// e.g. "ten thousand" → "ten thousandth", "one million" → "one millionth"
// The irregular map and tens map handle the non-regular forms.
var enLastWordOrdinal = map[string]string{
	// scale words
	"quintillion": "quintillionth",
	"quadrillion": "quadrillionth",
	"trillion":    "trillionth",
	"billion":     "billionth",
	"million":     "millionth",
	"thousand":    "thousandth",
	"hundred":     "hundredth",
	// tens
	"twenty":  "twentieth",
	"thirty":  "thirtieth",
	"forty":   "fortieth",
	"fifty":   "fiftieth",
	"sixty":   "sixtieth",
	"seventy": "seventieth",
	"eighty":  "eightieth",
	"ninety":  "ninetieth",
	// units/teens that need irregular handling
	"zero":      "zeroth",
	"one":       "first",
	"two":       "second",
	"three":     "third",
	"four":      "fourth",
	"five":      "fifth",
	"six":       "sixth",
	"seven":     "seventh",
	"eight":     "eighth",
	"nine":      "ninth",
	"ten":       "tenth",
	"eleven":    "eleventh",
	"twelve":    "twelfth",
	"thirteen":  "thirteenth",
	"fourteen":  "fourteenth",
	"fifteen":   "fifteenth",
	"sixteen":   "sixteenth",
	"seventeen": "seventeenth",
	"eighteen":  "eighteenth",
	"nineteen":  "nineteenth",
}

func (englishConverter) ToOrdinalWords(n int64) string {
	if n < 0 {
		// Negative ordinals are not meaningful for speech; prefix "minus"
		// to match Humanizer behaviour.
		return "minus " + englishOrdinal(-n)
	}
	return englishOrdinal(n)
}

// englishOrdinal converts n to its English ordinal form.
//
// Humanizer rule: when the full cardinal form starts with "one " (i.e. the
// number is an exact power-of-scale with multiplier 1, like 100, 1000,
// 1000000, or a compound like 100000 = "one hundred thousand"), the leading
// "one " is stripped before forming the ordinal.
// e.g. "one hundred" → "hundredth", "one hundred thousand" → "hundred thousandth"
// but "ten thousand" → "ten thousandth" (doesn't start with "one ")
// and "ten thousand one hundred and twenty-one" → "ten thousand one hundred
// and twenty-first" (internal "one" is not stripped).
func englishOrdinal(n int64) string {
	// Small values resolved directly via lookup tables.
	if w, ok := enOrdinalsIrregular[n]; ok {
		return w
	}
	if n < 20 {
		return enUnits[n] + "th"
	}
	if n < 100 {
		if w, ok := enTensOrdinal[n]; ok {
			return w
		}
		// Compound like twenty-first: ordinalise the unit part.
		return enTens[n/10] + "-" + englishOrdinal(n%10)
	}

	// Build the full cardinal, then apply ordinal transformation.
	cardinal := englishWords(n)

	// Strip a leading "one " only when the cardinal starts exactly with "one "
	// (meaning the number begins with the scale multiplier 1).
	body := cardinal
	if strings.HasPrefix(body, "one ") {
		body = body[4:] // strip "one "
	}

	// Make the last word ordinal. The last word may be hyphenated (e.g.
	// "twenty-one") — ordinalise only the part after the final hyphen.
	return ordinaliseLastWord(body)
}

// ordinaliseLastWord replaces the final word token in phrase with its ordinal
// form. Handles hyphenated compounds by ordinalising only the suffix after the
// last hyphen (e.g. "twenty-one" → "twenty-first").
func ordinaliseLastWord(phrase string) string {
	// Find the last space — everything after it is the final "word" (may
	// include a hyphen for compound tens).
	lastSpace := strings.LastIndex(phrase, " ")
	var prefix, lastWord string
	if lastSpace == -1 {
		prefix = ""
		lastWord = phrase
	} else {
		prefix = phrase[:lastSpace+1] // includes trailing space
		lastWord = phrase[lastSpace+1:]
	}

	// Handle hyphenated compound like "twenty-one": ordinalise only the
	// part after the last hyphen; keep the part before.
	hyphen := strings.LastIndex(lastWord, "-")
	var wordPrefix, stem string
	if hyphen != -1 {
		wordPrefix = lastWord[:hyphen+1] // e.g. "twenty-"
		stem = lastWord[hyphen+1:]       // e.g. "one"
	} else {
		wordPrefix = ""
		stem = lastWord
	}

	if ord, ok := enLastWordOrdinal[stem]; ok {
		return prefix + wordPrefix + ord
	}
	// Fallback: append "th" (covers edge cases not in the map).
	return prefix + wordPrefix + stem + "th"
}
