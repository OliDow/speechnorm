package speechnorm

import "strings"

func init() {
	Register("ar", arabicConverter{})
}

type arabicConverter struct{}

// arOnes: masculine cardinal units 0–19.
// Index 0 ("صفر") is used only by ToWords(0); the algorithm never indexes into
// arOnes[0] during group decomposition (groups with value 0 are skipped).
// Source: Humanizer ar.yml onesGroup.
var arOnes = [20]string{
	"صفر", "واحد", "اثنان", "ثلاثة", "أربعة",
	"خمسة", "ستة", "سبعة", "ثمانية", "تسعة",
	"عشرة", "أحد عشر", "اثنا عشر", "ثلاثة عشر", "أربعة عشر",
	"خمسة عشر", "ستة عشر", "سبعة عشر", "ثمانية عشر", "تسعة عشر",
}

// arTens: tens words indexed by tens digit (2–9 → 20–90).
// Indices 0 and 1 are unused; values 10–19 are handled by arOnes.
// Source: Humanizer ar.yml tensGroup.
var arTens = [10]string{
	"", "", "عشرون", "ثلاثون", "أربعون",
	"خمسون", "ستون", "سبعون", "ثمانون", "تسعون",
}

// arHundreds: hundreds words indexed by hundreds digit (1–9).
// Index 0 is unused. 200 (index 2) is the dual form "مئتان".
// 300–900 use the feminine construct-state multiplier + " مئة" because
// مئة is grammatically feminine — this is standard Arabic agreement,
// not a gender parameter choice.
// Source: Humanizer ar.yml hundredsGroup.
var arHundreds = [10]string{
	"", "مئة", "مئتان", "ثلاث مئة", "أربع مئة",
	"خمس مئة", "ست مئة", "سبع مئة", "ثمان مئة", "تسع مئة",
}

// arScaleWords holds the four forms needed for each scale level.
type arScaleWords struct {
	singular   string // value == 1 (appended, no prefix), or value%100 == 1, or terminal value > 10
	dual       string // value == 2 (appended, no prefix)
	plural     string // value 3–10
	accusative string // value > 10, non-terminal (tanween ـاً)
}

// arScales: scale words indexed by triad level.
// Level 0 (ones) has no scale word; levels 1–7 cover ألف through سكستيليون.
// Source: Humanizer ar.yml groups, appendedGroups, pluralGroups, appendedTwos.
var arScales = [8]arScaleWords{
	{},
	{"ألف", "ألفان", "آلاف", "ألفاً"},
	{"مليون", "مليونان", "ملايين", "مليوناً"},
	{"مليار", "ملياران", "مليارات", "ملياراً"},
	{"تريليون", "تريليونان", "تريليونات", "تريليوناً"},
	{"كوادريليون", "كوادريليونان", "كوادريليونات", "كوادريليوناً"},
	{"كوينتليون", "كوينتليونان", "كوينتليونات", "كوينتليوناً"},
	{"سكستيليون", "سكستيليونان", "سكستيليونات", "سكستيليوناً"},
}

func (arabicConverter) ToWords(n int64) string {
	if n == 0 {
		return "صفر"
	}
	if n < 0 {
		return "سالب " + arabicWords(-n)
	}
	return arabicWords(n)
}

func (arabicConverter) ToOrdinalWords(n int64) string {
	return arabicConverter{}.ToWords(n)
}

// arabicWords converts n > 0 to Arabic cardinal words.
func arabicWords(n int64) string {
	// Decompose into triads (groups of 3 digits), least significant first.
	type triad struct {
		value int64
		level int
	}
	var triads []triad
	remaining := n
	for level := 0; remaining > 0; level++ {
		triads = append(triads, triad{value: remaining % 1000, level: level})
		remaining /= 1000
	}

	// Process from most significant to least significant.
	var parts []string
	for i := len(triads) - 1; i >= 0; i-- {
		t := triads[i]
		if t.value == 0 {
			continue
		}

		if t.level == 0 {
			parts = append(parts, arabicUnder1000(t.value))
			continue
		}

		// Check whether any lower-level group is non-zero. This determines
		// whether the scale word takes the accusative (tanween) form.
		hasLower := false
		for j := i - 1; j >= 0; j-- {
			if triads[j].value > 0 {
				hasLower = true
				break
			}
		}

		parts = append(parts, arabicScaleGroup(t.value, t.level, hasLower))
	}

	return strings.Join(parts, " و ")
}

// arabicScaleGroup renders a single triad at the given scale level
// (thousands, millions, etc.) including its scale word.
func arabicScaleGroup(value int64, level int, hasLowerGroups bool) string {
	scale := arScales[level]

	// Appended forms: value 1 or 2 use the scale word alone with no
	// numeric prefix ("ألف" not "واحد ألف", "ألفان" not "اثنان ألف").
	if value == 1 {
		return scale.singular
	}
	if value == 2 {
		return scale.dual
	}

	numberPart := arabicUnder1000(value)

	// When the group value ends in 1 (e.g. 101, 201), use the singular
	// scale form regardless of terminal/non-terminal position.
	// Source: Humanizer AppendedGroupNumberToWordsConverter.cs,
	// `groupNumber % 100 != 1` branch.
	if value%100 == 1 {
		return numberPart + " " + scale.singular
	}

	// 3–10: plural form.
	if value >= 3 && value <= 10 {
		return numberPart + " " + scale.plural
	}

	// > 10: accusative (tanween) when non-terminal, singular when terminal.
	if hasLowerGroups {
		return numberPart + " " + scale.accusative
	}
	return numberPart + " " + scale.singular
}

// arabicUnder1000 converts 1–999 to Arabic words.
func arabicUnder1000(n int64) string {
	hundreds := n / 100
	remainder := n % 100

	// Special case: exactly 200 uses the dual form.
	if hundreds == 2 && remainder == 0 {
		return "مئتان"
	}

	var parts []string

	if hundreds > 0 {
		parts = append(parts, arHundreds[hundreds])
	}

	if remainder > 0 {
		parts = append(parts, arabicUnder100(remainder))
	}

	return strings.Join(parts, " و ")
}

// arabicUnder100 converts 1–99 to Arabic words.
func arabicUnder100(n int64) string {
	if n < 20 {
		return arOnes[n]
	}

	tensDigit := n / 10
	unit := n % 10

	if unit == 0 {
		return arTens[tensDigit]
	}

	// Arabic: units before tens, joined with و (like German's units-before-tens).
	return arOnes[unit] + " و " + arTens[tensDigit]
}
