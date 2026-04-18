package speechnorm

import "testing"

// TestArabicToWords is a verbatim port of every ar cardinal masculine row from
// Humanizer's LocaleNumberTheoryData.cs and LocaleAdditionalNumberTheoryData.cs.
// Do not add, remove, or modify cases without regenerating from Humanizer.
//
// Skipped: GrammaticalGender.Feminine rows — speechnorm uses masculine default;
//
//	see TestArabicToWords_FeminineSkipped.
//
// Skipped: ordinal rows — see TestArabicToOrdinalWords_Skipped.
// Skipped: WordForm rows — output matches default cardinal for Arabic.
// Skipped: CardinalAddAndCases — output matches default cardinal for Arabic
//
//	(addAnd flag has no effect on Arabic output).
//
// Humanizer correction: some large-number test cases in
// LocaleNumberOverloadTheoryData.cs contain extra whitespace artifacts from the
// C# engine's handling of consecutive appended groups (e.g., "مليون و   ألف").
// Our implementation produces clean single-space output; those values are tested
// in numwords_ar_scale_test.go with normalised whitespace and a comment noting
// the Humanizer discrepancy.
func TestArabicToWords(t *testing.T) {
	conv, ok := lookup("ar")
	if !ok {
		t.Fatal("arabic converter not registered")
	}
	cases := []struct {
		n    int64
		want string
	}{
		// --- LocaleNumberTheoryData.cs CardinalCases ---
		{0, "صفر"},
		{1, "واحد"},
		{2, "اثنان"},
		{3, "ثلاثة"},
		{4, "أربعة"},
		{5, "خمسة"},
		{10, "عشرة"},
		{11, "أحد عشر"},
		{12, "اثنا عشر"},
		{19, "تسعة عشر"},
		{20, "عشرون"},
		{21, "واحد و عشرون"},
		{22, "اثنان و عشرون"},
		{30, "ثلاثون"},
		{40, "أربعون"},
		{80, "ثمانون"},
		{90, "تسعون"},
		{99, "تسعة و تسعون"},
		{100, "مئة"},
		{101, "مئة و واحد"},
		{110, "مئة و عشرة"},
		{111, "مئة و أحد عشر"},
		{115, "مئة و خمسة عشر"},
		{121, "مئة و واحد و عشرون"},
		{200, "مئتان"},
		{999, "تسع مئة و تسعة و تسعون"},
		{1000, "ألف"},
		{1001, "ألف و واحد"},

		// --- LocaleAdditionalNumberTheoryData.cs AdditionalCardinalCases (masculine) ---
		{-1, "سالب واحد"},
		{-2, "سالب اثنان"},
		{-11, "سالب أحد عشر"},
		{-22, "سالب اثنان و عشرون"},
		{1111, "ألف و مئة و أحد عشر"},
		{3501, "ثلاثة آلاف و خمس مئة و واحد"},
		{-3501, "سالب ثلاثة آلاف و خمس مئة و واحد"},
		{11111, "أحد عشر ألفاً و مئة و أحد عشر"},
		{111111, "مئة و أحد عشر ألفاً و مئة و أحد عشر"},
		{1000001, "مليون و واحد"},
		{-1000001, "سالب مليون و واحد"},
		{1111111, "مليون و مئة و أحد عشر ألفاً و مئة و أحد عشر"},
		{11111111, "أحد عشر مليوناً و مئة و أحد عشر ألفاً و مئة و أحد عشر"},
		{111111111, "مئة و أحد عشر مليوناً و مئة و أحد عشر ألفاً و مئة و أحد عشر"},
		{1111111111, "مليار و مئة و أحد عشر مليوناً و مئة و أحد عشر ألفاً و مئة و أحد عشر"},
		{11111111111, "أحد عشر ملياراً و مئة و أحد عشر مليوناً و مئة و أحد عشر ألفاً و مئة و أحد عشر"},
		{111111111111, "مئة و أحد عشر ملياراً و مئة و أحد عشر مليوناً و مئة و أحد عشر ألفاً و مئة و أحد عشر"},
		{1111111111111, "تريليون و مئة و أحد عشر ملياراً و مئة و أحد عشر مليوناً و مئة و أحد عشر ألفاً و مئة و أحد عشر"},
		{11111111111111, "أحد عشر تريليوناً و مئة و أحد عشر ملياراً و مئة و أحد عشر مليوناً و مئة و أحد عشر ألفاً و مئة و أحد عشر"},
		{111111111111111, "مئة و أحد عشر تريليوناً و مئة و أحد عشر ملياراً و مئة و أحد عشر مليوناً و مئة و أحد عشر ألفاً و مئة و أحد عشر"},
		{1111111111111111, "كوادريليون و مئة و أحد عشر تريليوناً و مئة و أحد عشر ملياراً و مئة و أحد عشر مليوناً و مئة و أحد عشر ألفاً و مئة و أحد عشر"},
		{11111111111111111, "أحد عشر كوادريليوناً و مئة و أحد عشر تريليوناً و مئة و أحد عشر ملياراً و مئة و أحد عشر مليوناً و مئة و أحد عشر ألفاً و مئة و أحد عشر"},
		{111111111111111111, "مئة و أحد عشر كوادريليوناً و مئة و أحد عشر تريليوناً و مئة و أحد عشر ملياراً و مئة و أحد عشر مليوناً و مئة و أحد عشر ألفاً و مئة و أحد عشر"},
		{1111111111111111111, "كوينتليون و مئة و أحد عشر كوادريليوناً و مئة و أحد عشر تريليوناً و مئة و أحد عشر ملياراً و مئة و أحد عشر مليوناً و مئة و أحد عشر ألفاً و مئة و أحد عشر"},
		{10000000001, "عشرة مليارات و واحد"},
		{-10000000001, "سالب عشرة مليارات و واحد"},
		{8750000500001, "ثمانية تريليونات و سبع مئة و خمسون ملياراً و خمس مئة ألفاً و واحد"},
		{-8750000500001, "سالب ثمانية تريليونات و سبع مئة و خمسون ملياراً و خمس مئة ألفاً و واحد"},
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

// TestArabicToOrdinalWords_Skipped documents every ar ordinal row from
// Humanizer's LocaleAdditionalNumberTheoryData.cs. These are skipped
// because speechnorm's ordinal regex is English-suffix-specific (st|nd|rd|th)
// and ToOrdinalWords returns cardinal form for all non-English locales.
// The rows are preserved here for completeness and future reference.
func TestArabicToOrdinalWords_Skipped(t *testing.T) {
	t.Skip("speechnorm ordinal regex is English-only; ToOrdinalWords returns cardinal for ar")

	// Masculine ordinal rows from LocaleAdditionalNumberTheoryData.cs:
	// {0, "الصفر"}, {1, "الأول"}, {2, "الثاني"}, {3, "الثالث"},
	// {4, "الرابع"}, {5, "الخامس"}, {6, "السادس"}, {7, "السابع"},
	// {8, "الثامن"}, {9, "التاسع"}, {10, "العاشر"},
	// {11, "الحادي عشر"}, {12, "الثاني عشر"}, {13, "الثالث عشر"},
	// {14, "الرابع عشر"}, {15, "الخامس عشر"}, {16, "السادس عشر"},
	// {17, "السابع عشر"}, {18, "الثامن عشر"}, {19, "التاسع عشر"},
	// {20, "العشرون"}, {21, "الحادي و العشرون"}, {22, "الثاني و العشرون"},
	// {30, "الثلاثون"}, {40, "الأربعون"}, {50, "الخمسون"},
	// {60, "الستون"}, {70, "السبعون"}, {80, "الثمانون"}, {90, "التسعون"},
	// {95, "الخامس و التسعون"}, {96, "السادس و التسعون"},
	// {100, "المئة"}, {120, "العشرون بعد المئة"},
	// {121, "الحادي و العشرون بعد المئة"},
	// {200, "المئتان"}, {221, "الحادي و العشرون بعد المئتان"},
	// {300, "الثلاث مئة"}, {321, "الحادي و العشرون بعد الثلاث مئة"},
	// {327, "السابع و العشرون بعد الثلاث مئة"},
	// {1000, "الألف"}, {1001, "الأول بعد الألف"},
	// {1021, "الحادي و العشرون بعد الألف"},
	// {10000, "العشرة آلاف"}, {10121, "الحادي و العشرون بعد العشرة آلاف و مئة"},
	// {100000, "المئة ألف"}, {1000000, "المليون"},
	// {1020135, "الخامس و الثلاثون بعد المليون و عشرون ألفاً و مئة"},
	//
	// Feminine ordinal rows (GrammaticalGender.Feminine):
	// {0, "الصفر"}, {1, "الأولى"}, {2, "الثانية"}, {3, "الثالثة"},
	// {4, "الرابعة"}, {5, "الخامسة"}, {6, "السادسة"}, {7, "السابعة"},
	// {8, "الثامنة"}, {9, "التاسعة"}, {10, "العاشرة"},
	// {11, "الحادية عشرة"}, {12, "الثانية عشرة"}, {13, "الثالثة عشرة"},
	// {14, "الرابعة عشرة"}, {15, "الخامسة عشرة"}, {16, "السادسة عشرة"},
	// {17, "السابعة عشرة"}, {18, "الثامنة عشرة"}, {19, "التاسعة عشرة"},
	// {20, "العشرون"}, {21, "الحادية و العشرون"}, {22, "الثانية و العشرون"},
	// {30, "الثلاثون"}, {40, "الأربعون"}, {50, "الخمسون"},
	// {60, "الستون"}, {70, "السبعون"}, {80, "الثمانون"}, {90, "التسعون"},
	// {95, "الخامسة و التسعون"}, {96, "السادسة و التسعون"},
	// {100, "المئة"}, {120, "العشرون بعد المئة"},
	// {121, "الحادية و العشرون بعد المئة"},
	// {200, "المئتان"}, {221, "الحادية و العشرون بعد المئتان"},
	// {300, "الثلاث مئة"}, {321, "الحادية و العشرون بعد الثلاث مئة"},
	// {327, "السابعة و العشرون بعد الثلاث مئة"},
	// {1000, "الألف"}, {1001, "الأولى بعد الألف"},
	// {1021, "الحادية و العشرون بعد الألف"},
	// {10000, "العشرة آلاف"}, {10121, "الحادية و العشرون بعد العشرة آلاف و مئة"},
	// {100000, "المئة ألف"}, {1000000, "المليون"},
	// {1020135, "الخامسة و الثلاثون بعد المليون و عشرون ألفاً و مئة"},
	//
	// Also from LocaleNumberTheoryData.cs OrdinalCases:
	// {1, "الأول"}, {2, "الثاني"}, {3, "الثالث"}, {4, "الرابع"},
	// {10, "العاشر"}, {11, "الحادي عشر"}, {12, "الثاني عشر"},
	// {20, "العشرون"}, {21, "الحادي و العشرون"},
	// {100, "المئة"}, {101, "الأول بعد المئة"},
}

// TestArabicToWords_FeminineSkipped documents every ar feminine cardinal row
// from Humanizer. These are skipped because speechnorm uses masculine default
// (the NormaliseNumbers API has no gender parameter).
func TestArabicToWords_FeminineSkipped(t *testing.T) {
	t.Skip("speechnorm uses masculine default; feminine form not supported")

	// Feminine cardinal rows from LocaleNumberTheoryData.cs CardinalGenderCases:
	// {1, GrammaticalGender.Feminine, "واحدة"},
	//
	// Feminine cardinal rows from LocaleAdditionalNumberTheoryData.cs:
	// {0, "صفر"}, {1, "واحدة"}, {2, "اثنتان"},
	// {11, "إحدى عشرة"}, {22, "اثنتان و عشرون"},
	// {122, "مئة و اثنتان و عشرون"},
	// {3501, "ثلاثة آلاف و خمس مئة و واحدة"},
	// {1000001, "مليون و واحدة"},
	// {10000000001, "عشرة مليارات و واحدة"},
	// {8750000500001, "ثمانية تريليونات و سبع مئة و خمسون ملياراً و خمس مئة ألفاً و واحدة"},
	//
	// CardinalWordFormGenderCases:
	// {21, Feminine, "واحدة و عشرون"},
}
