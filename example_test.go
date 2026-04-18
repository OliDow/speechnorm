package speechnorm_test

import (
	"fmt"

	"github.com/olidow/speechnorm"
)

func ExampleNormaliseNumbers() {
	fmt.Println(speechnorm.NormaliseNumbers("I have 3 cats and $5", "en"))
	// Output: I have three cats and five dollars
}

func ExampleNormaliseNumbers_german() {
	fmt.Println(speechnorm.NormaliseNumbers("Ich habe 3 Katzen", "de"))
	// Output: Ich habe drei Katzen
}

func ExampleNormaliseNumbers_spanish() {
	fmt.Println(speechnorm.NormaliseNumbers("Tengo 3 gatos", "es"))
	// Output: Tengo tres gatos
}

func ExampleNormaliseNumbers_french() {
	fmt.Println(speechnorm.NormaliseNumbers("J'ai 3 chats", "fr"))
	// Output: J'ai trois chats
}

func ExampleNormaliseNumbers_italian() {
	fmt.Println(speechnorm.NormaliseNumbers("Ho 3 gatti", "it"))
	// Output: Ho tre gatti
}

func ExampleNormaliseNumbers_portuguese() {
	fmt.Println(speechnorm.NormaliseNumbers("Tenho 3 gatos", "pt"))
	// Output: Tenho três gatos
}

func ExampleNormaliseNumbers_arabic() {
	fmt.Println(speechnorm.NormaliseNumbers("لدي 3 قطط", "ar"))
	// Output: لدي ثلاثة قطط
}
