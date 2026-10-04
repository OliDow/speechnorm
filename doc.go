// Package speechnorm rewrites digit patterns in free-form text into
// locale-appropriate spoken words for TTS input. Supported locales:
// ar, de, en, es, fr, it, pt. The single entry point is NormaliseNumbers.
//
// Currency amounts are spoken in the target locale: number words, joiners,
// and unit names all follow the locale's converter.
//
// The package has zero non-stdlib dependencies.
package speechnorm
