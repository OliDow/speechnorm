package speechnorm

import (
	"regexp"
	"strconv"
	"strings"
)

// NormaliseNumbers converts digit patterns in text (currency, ordinals,
// comma-grouped integers, plain integers) into locale-appropriate words.
// Returns text unchanged when locale is empty or unregistered. Pure
// function: no side effects, no logging.
func NormaliseNumbers(text, locale string) string {
	if text == "" || strings.TrimSpace(locale) == "" {
		return text
	}
	conv, ok := lookup(locale)
	if !ok {
		return text
	}

	text = currencyRegex.ReplaceAllStringFunc(text, func(m string) string {
		return replaceCurrency(m, conv)
	})
	text = ordinalRegex.ReplaceAllStringFunc(text, func(m string) string {
		return replaceOrdinal(m, conv)
	})
	text = commaNumberRegex.ReplaceAllStringFunc(text, func(m string) string {
		return replaceCommaNumber(m, conv)
	})
	text = replacePlainIntegers(text, conv)
	return text
}

var (
	// \p{Sc} = Unicode currency symbols
	currencyRegex    = regexp.MustCompile(`(\p{Sc})\s?(\d{1,3}(?:,\d{3})*)(?:\.(\d{1,2}))?`)
	ordinalRegex     = regexp.MustCompile(`\b(\d{1,19})(st|nd|rd|th)\b`)
	commaNumberRegex = regexp.MustCompile(`\b\d{1,3}(?:,\d{3})+\b`)
	plainIntRegex    = regexp.MustCompile(`\b\d{1,19}\b`)
)

func replaceCurrency(match string, conv Converter) string {
	sub := currencyRegex.FindStringSubmatch(match)
	if sub == nil {
		return match
	}
	symbol := sub[1]
	integerPart := sub[2]
	decimalPart := ""
	if len(sub) > 3 {
		decimalPart = sub[3]
	}
	return convertCurrency(symbol, integerPart, decimalPart, conv)
}

func replaceOrdinal(match string, conv Converter) string {
	sub := ordinalRegex.FindStringSubmatch(match)
	if sub == nil {
		return match
	}
	n, err := strconv.ParseInt(sub[1], 10, 64)
	if err != nil {
		return match
	}
	return conv.ToOrdinalWords(n)
}

func replaceCommaNumber(match string, conv Converter) string {
	stripped := strings.ReplaceAll(match, ",", "")
	n, err := strconv.ParseInt(stripped, 10, 64)
	if err != nil {
		return match
	}
	return conv.ToWords(n)
}

// replacePlainIntegers walks the text, finds \b\d{1,12}\b matches, and
// rejects those adjacent to '.' or ',' on either side — reproducing the
// C# negative lookbehind/lookahead behaviour without lookaround (which
// RE2 forbids).
func replacePlainIntegers(text string, conv Converter) string {
	var b strings.Builder
	indices := plainIntRegex.FindAllStringIndex(text, -1)
	cursor := 0
	for _, idx := range indices {
		start, end := idx[0], idx[1]
		// Emit anything before this match unchanged.
		b.WriteString(text[cursor:start])

		// Preceding-byte check.
		if start > 0 {
			prev := text[start-1]
			if prev == '.' || prev == ',' {
				b.WriteString(text[start:end])
				cursor = end
				continue
			}
		}
		// Trailing-byte check: '.' or ',' followed by a digit means we're
		// looking at a decimal, leave unchanged.
		if end < len(text) {
			next := text[end]
			if (next == '.' || next == ',') && end+1 < len(text) {
				after := text[end+1]
				if after >= '0' && after <= '9' {
					b.WriteString(text[start:end])
					cursor = end
					continue
				}
			}
		}

		// Safe to convert.
		n, err := strconv.ParseInt(text[start:end], 10, 64)
		if err != nil {
			b.WriteString(text[start:end])
		} else {
			b.WriteString(conv.ToWords(n))
		}
		cursor = end
	}
	b.WriteString(text[cursor:])
	return b.String()
}
