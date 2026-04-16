// Package speechnorm rewrites digit patterns in free-form text into
// locale-appropriate spoken words for TTS input. Supported locales:
// en, de, es, fr, it, pt. The single entry point is NormaliseNumbers.
//
// Currency words ("dollars", "and fifty cents") are always English
// regardless of locale; only the number words follow the locale's
// converter.
//
// The package has zero non-stdlib dependencies.
package speechnorm
