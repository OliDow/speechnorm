package speechnorm

import "testing"

// TestArabicToWords_ScaleBoundaries covers scale values not in
// Humanizer's main LocaleNumberTheoryData.cs / AdditionalCardinalCases.
//
// Sources:
// - LocaleNumberMagnitudeTheoryData.cs (ar rows, non-duplicate)
// - Generated scale-boundary cases at each power-of-ten transition
//
// Humanizer whitespace note: the C# engine produces extra whitespace
// for consecutive appended groups (e.g., "مليون و   ألف" for 1001001).
// Our implementation produces clean single-space output. Expected strings
// here use normalised whitespace. See AppendedGroupNumberToWordsConverter.cs
// tens==1/groupLevel>0/hundreds==0 branch for the source of the C# artifact.
func TestArabicToWords_ScaleBoundaries(t *testing.T) {
	conv, _ := lookup("ar")
	cases := []struct {
		n    int64
		want string
	}{
		// Dual forms
		{2000, "ألفان"},
		{2000000, "مليونان"},
		{2000000000, "ملياران"},

		// Plural range (3–10)
		{3000, "ثلاثة آلاف"},
		{5000, "خمسة آلاف"},
		{10000, "عشرة آلاف"},
		{3000000, "ثلاثة ملايين"},
		{10000000, "عشرة ملايين"},
		{3000000000, "ثلاثة مليارات"},
		{10000000000, "عشرة مليارات"},

		// Post-plural singular (> 10, terminal — no tanween)
		{100000, "مئة ألف"},
		{100000000, "مئة مليون"},

		// Compound values crossing scale boundaries
		{1100, "ألف و مئة"},
		{2001, "ألفان و واحد"},
		{1000000000000, "تريليون"},
		{1000000000000000, "كوادريليون"},

		// --- LocaleNumberMagnitudeTheoryData.cs (normalised whitespace) ---
		{1001000001, "مليار و مليون و واحد"},
		{4325010007018, "أربعة تريليونات و ثلاث مئة و خمسة و عشرون ملياراً و عشرة ملايين و سبعة آلاف و ثمانية عشر"},
		{1000000000000000000, "كوينتليون"},
		{1000000000000000001, "كوينتليون و واحد"},
	}
	for _, c := range cases {
		c := c
		t.Run(intToBase10(c.n), func(t *testing.T) {
			if got := conv.ToWords(c.n); got != c.want {
				t.Errorf("ar ToWords(%d) = %q, want %q", c.n, got, c.want)
			}
		})
	}
}
