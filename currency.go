package speechnorm

import (
	"strconv"
	"strings"
)

// currencyNames returns (majorSingular, majorPlural, minorSingular, minorPlural)
// for a currency symbol. All four return empty for unknown symbols.
func currencyNames(symbol string) (string, string, string, string) {
	switch symbol {
	case "$":
		return "dollar", "dollars", "cent", "cents"
	case "€":
		return "euro", "euros", "cent", "cents"
	case "£":
		return "pound", "pounds", "penny", "pence"
	case "¥":
		return "yen", "yen", "", ""
	case "₹":
		return "rupee", "rupees", "paisa", "paise"
	}
	return "", "", "", ""
}

// convertCurrency formats a currency amount into spoken words. The input
// integerPart may contain commas (they are stripped). decimalPart is 0-2
// digits; a single digit is padded to two (so $5.5 reads as "fifty cents").
// The caller is responsible for supplying digits-only content; we assume
// regex-validated input.
func convertCurrency(symbol, integerPart, decimalPart string, conv Converter) string {
	major, _ := strconv.ParseInt(strings.ReplaceAll(integerPart, ",", ""), 10, 64)

	majorSing, majorPlu, minorSing, minorPlu := currencyNames(symbol)

	if majorSing == "" {
		return conv.ToWords(major)
	}

	majorUnit := majorPlu
	if major == 1 {
		majorUnit = majorSing
	}
	result := conv.ToWords(major) + " " + majorUnit

	if decimalPart == "" {
		return result
	}

	dec := decimalPart
	if len(dec) == 1 {
		dec += "0"
	}

	minor, _ := strconv.ParseInt(dec, 10, 64)
	if minor == 0 || minorSing == "" {
		return result
	}

	minorUnit := minorPlu
	if minor == 1 {
		minorUnit = minorSing
	}
	return result + " and " + conv.ToWords(minor) + " " + minorUnit
}
